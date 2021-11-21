package skillz

type Operator struct {
	Column    string    `json:"column"`
	Operation Operation `json:"operation"`
	Value     OpValue   `json:"value"`
}

type OpValue interface{}

type Operation int

const (
	EqualOp Operation = iota
	NotEqualOp
	GreaterThanOp
	GreaterThanEqualToOp
	LessThanOp
	LessThanEqualToOp
	InOp
	NotInOp
	LikeOp

	LimitOp
	OrderOp
	SkipOp
	OrOp
	AndOp
	ExistsOp
)

var AllOperations = []Operation{
	EqualOp,
	NotEqualOp,
	GreaterThanOp,
	GreaterThanEqualToOp,
	LessThanOp,
	LessThanEqualToOp,
	InOp,
	NotInOp,
	LikeOp,
	LimitOp,
	OrderOp,
	SkipOp,
	OrOp,
	AndOp,
	ExistsOp,
}

func (o Operation) IsValid() bool {
	switch o {
	case EqualOp, NotEqualOp,
		GreaterThanOp, LessThanOp, GreaterThanEqualToOp, LessThanEqualToOp,
		InOp, NotInOp, LikeOp,
		LimitOp, OrderOp, SkipOp, OrOp, AndOp, ExistsOp:
		return true
	}
	return false
}

func NewEqualOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: EqualOp,
		Value:     value,
	}
}

func NewNotEqualOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: NotEqualOp,
		Value:     value,
	}
}

func NewGreaterThanOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: GreaterThanOp,
		Value:     value,
	}
}

func NewGreaterThanEqualToOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: GreaterThanEqualToOp,
		Value:     value,
	}
}

func NewLessThanOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: LessThanOp,
		Value:     value,
	}
}

func NewLessThanEqualToOperator(column string, value interface{}) *Operator {
	return &Operator{
		Column:    column,
		Operation: LessThanEqualToOp,
		Value:     value,
	}
}

type Sort int

const (
	SortAsc  Sort = 1
	SortDesc Sort = -1
)

var AllSort = []Sort{
	SortAsc,
	SortDesc,
}

func (e Sort) IsValid() bool {
	switch e {
	case SortAsc, SortDesc:
		return true
	}
	return false
}

func (e Sort) Value() int {
	return int(e)
}

func NewOrderOperator(column string, sort Sort) *Operator {

	if !sort.IsValid() {
		return nil
	}

	return &Operator{
		Column:    column,
		Operation: OrderOp,
		Value:     sort.Value(),
	}

}

func NewInOperator(column string, value []interface{}) *Operator {

	return &Operator{
		Column:    column,
		Operation: InOp,
		Value:     value,
	}

}

func NewNotInOperator(column string, value interface{}) *Operator {

	return &Operator{
		Column:    column,
		Operation: NotInOp,
		Value:     value,
	}

}

func NewLimitOperator(value uint64) *Operator {
	return &Operator{
		Column:    "",
		Operation: LimitOp,
		Value:     value,
	}
}

func NewSkipOperator(value int64) *Operator {
	return &Operator{
		Column:    "",
		Operation: SkipOp,
		Value:     value,
	}
}

func NewOrOperator(value ...*Operator) *Operator {
	return &Operator{
		Column:    "",
		Operation: OrOp,
		Value:     value,
	}
}

func NewAndOperator(value ...*Operator) *Operator {
	return &Operator{
		Column:    "",
		Operation: AndOp,
		Value:     value,
	}
}

func NewExistsOperator(column string, value bool) *Operator {
	return &Operator{
		Column:    column,
		Operation: ExistsOp,
		Value:     value,
	}
}

func NewLikeOperator(column string, value string) *Operator {
	return &Operator{
		Column:    column,
		Operation: LikeOp,
		Value:     value,
	}
}
