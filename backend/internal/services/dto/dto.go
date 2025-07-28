package dto

// 通用 Model -> DTO 轉換

// 泛型包裝：統一呼叫方式
func ModelToDTO[T any, D any](model *T, convert func(*T) *D) *D {
	return convert(model)
}

// 可依需求擴充其它 DTO 轉換
