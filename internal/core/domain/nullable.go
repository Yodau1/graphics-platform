package domain

// Nullable — обобщённый тип для работы с «трёхзначной» логикой
// Значения полей:
//   - Set=false           → поле не было передано (Value не имеет смысла) → ничего не делаем
//   - Set=true, Value=nil → поле передано со значением *nil               → например, очистить в БД
//   - Set=true, Value=&v  → поле передано с конкретным значением v        → например, обновить в БД

type Nullable[T any] struct {
	Value *T
	Set   bool
}
