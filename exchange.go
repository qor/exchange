package exchange

import (
	"errors"
	"fmt"
	"reflect"

	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/roles"
	"github.com/qor/validations"
)

// Resource defined an exchange resource, which includes importing/exporting fields definitions
type Resource struct {
	resource.Resource
	Config *Config
	Metas  []*Meta
}

// Config is exchange resource config
type Config struct {
	// PrimaryField that used as primary field when searching resource from database
	PrimaryField string
	// Permission defined permission
	Permission *roles.Permission
	// WithoutHeader no header in the data file
	WithoutHeader bool
}

// NewResource new exchange Resource
func NewResource(value interface{}, config ...Config) *Resource {
	res := Resource{Resource: *resource.New(value)}
	if len(config) > 0 {
		res.Config = &config[0]
	} else {
		res.Config = &Config{}
	}
	res.Permission = res.Config.Permission

	if res.Config.PrimaryField != "" {
		res.FindOneHandler = func(result interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
			scope := context.GetDB().NewScope(res.Value)
			field, ok := scope.FieldByName(res.Config.PrimaryField)
			if !ok {
				return errors.New("failed to find primary field")
			}
			primaryCond := fmt.Sprintf("%v = ?", scope.Quote(field.DBName))
			primaryMeta := metaValues.Get(res.Config.PrimaryField)
			if primaryMeta == nil {
				return fmt.Errorf("primary field \"%s\" not exist in meta values: %s", field.Name, jsonV(metaValues))
			}
			primaryValue := metaValues.Get(res.Config.PrimaryField).Value
			return context.GetDB().First(result, primaryCond, primaryValue).Error
		}
	}
	return &res
}

func jsonV(v interface{}) (r string) {
	b, _ := json.MarshalIndent(v, "", "\t")
	r = string(b)
	return
}

// Meta define exporting/importing meta for exchange Resource
func (res *Resource) Meta(meta *Meta) *Meta {
	if meta.Header == "" {
		meta.Header = meta.Name
	}

	meta.base = res
	meta.updateMeta()
	res.Metas = append(res.Metas, meta)
	return meta
}

// GetMeta get defined Meta from exchange Resource
func (res *Resource) GetMeta(name string) *Meta {
	for _, meta := range res.Metas {
		if meta.Header == name {
			return meta
		}
	}
	return nil
}

// GetMetas get all defined Metas from exchange Resource
func (res *Resource) GetMetas([]string) []resource.Metaor {
	metas := []resource.Metaor{}
	for _, meta := range res.Metas {
		metas = append(metas, meta)
	}
	return metas
}

type errorsInterface interface {
	GetErrors() []error
}

// Import used to import data into a exchange Resource
//     product.Import(csv.New("products.csv"), context)
func (res *Resource) Import(container Container, context *qor.Context, callbacks ...func(Progress) error) (err error) {
	var rows Rows
	rows, err = container.NewReader(res, context)
	if err != nil {
		return err
	}

	var current uint
	var total = rows.Total()

	if db := context.GetDB(); db != nil {
		tx := db.Begin()
		context.SetDB(tx)
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}

	for rows.Next() {
		current++
		progress := Progress{Total: total, Current: current}

		var metaValues *resource.MetaValues
		metaValues, err = rows.ReadRow()
		if err != nil {
			handleError(err, &progress)
			continue
		}

		for _, metaValue := range metaValues.Values {
			progress.Cells = append(progress.Cells, Cell{
				Header: metaValue.Name,
				Value:  metaValue.Value,
			})
		}

		result := res.NewStruct()
		progress.Value = result

		err = res.FindOneHandler(result, metaValues, context)
		if err != nil && err != gorm.ErrRecordNotFound {
			handleError(err, &progress)
			continue
		}

		err = resource.DecodeToResource(res, result, metaValues, context).Start()
		if err != nil {
			handleError(err, &progress)
			continue
		}

		err = res.CallSave(result, context)
		if err != nil {
			handleError(err, &progress)
			continue
		}

		for _, callback := range callbacks {
			if err = callback(progress); err != nil {
				handleError(err, &progress)
				continue
			}
		}
	}
	return err
}

func handleError(err error, progress *Progress) {
	if errors, ok := err.(errorsInterface); ok {
		for _, err := range errors.GetErrors() {
			handleError(err, progress)
		}
	} else if err, ok := err.(*validations.Error); ok {
		for idx, cell := range progress.Cells {
			if cell.Header == err.Column {
				cell.Error = err
				progress.Cells[idx] = cell
				break
			}
		}
	} else if len(progress.Cells) > 0 {
		var err error = err
		cell := progress.Cells[0]
		if cell.Error != nil {
			var errors qor.Errors
			errors.AddError(cell.Error)
			errors.AddError(err)
			err = errors
		}
		cell.Error = err
	}
}

// Export used export data from a exchange Resource
//     product.Export(csv.New("products.csv"), context)
func (res *Resource) Export(container Container, context *qor.Context, callbacks ...func(Progress) error) error {
	results := res.NewSlice()

	var total uint
	if err := context.GetDB().Find(results).Count(&total).Error; err == nil {
		reflectValue := reflect.Indirect(reflect.ValueOf(results))

		if writer, err := container.NewWriter(res, context); err == nil {
			writer.WriteHeader()

			for i := 0; i < reflectValue.Len(); i++ {
				var result = reflectValue.Index(i).Interface()
				var metaValues *resource.MetaValues
				if metaValues, err = writer.WriteRow(result); err != nil {
					return err
				}

				var progress = Progress{
					Current: uint(i + 1),
					Total:   total,
					Value:   result,
				}

				for _, metaValue := range metaValues.Values {
					progress.Cells = append(progress.Cells, Cell{
						Header: metaValue.Name,
						Value:  metaValue.Value,
					})
				}

				for _, callback := range callbacks {
					if err := callback(progress); err != nil {
						return err
					}
				}
			}
			writer.Flush()
		} else {
			return err
		}
	} else {
		return err
	}
	return nil
}
