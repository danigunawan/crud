package crud

type RowValue struct {
	SQLColumn string
	Value     interface{}
}

type Row struct {
	SQLTableName string
	Values       []*RowValue
}

func (row *Row) SQLValues() map[string]interface{} {
	result := map[string]interface{}{}

	for _, v := range row.Values {
		result[v.SQLColumn] = v.Value
	}

	return result
}

func NewRow(st interface{}) (*Row, error) {
	values, err := GetRowValuesOf(st)
	if err != nil {
		return nil, err
	}

	return &Row{
		SQLTableName: SQLTableNameOf(st),
		Values:       values,
	}, nil
}

func GetRowValuesOf(st interface{}) ([]*RowValue, error) {
	fields, err := CollectRows(st, []*RowValue{})
	if err != nil {
		return nil, err
	}

	return fields, nil
}

func CollectRows(st interface{}, rows []*RowValue) ([]*RowValue, error) {
	iter := NewFieldIteration(st)
	for iter.Next() {
		if iter.IsEmbeddedStruct() {
			if _rows, err := CollectRows(iter.ValueField().Interface(), rows); err != nil {
				return nil, err
			} else {
				rows = _rows
			}
			continue
		}

		sqlOptions, err := iter.SQLOptions()

		if err != nil {
			return nil, err
		}

		if sqlOptions.Ignore {
			continue
		}

		value := iter.Value()

		if n, ok := value.(int); ok && sqlOptions.AutoIncrement > 0 && n == 0 {
			continue
		}

		rows = append(rows, &RowValue{
			SQLColumn: sqlOptions.Name,
			Value:     value,
		})
	}

	return rows, nil
}
