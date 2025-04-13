package main

func int64ToInt32Ptr(i int64) int32 {
	i32 := int32(i)
	return i32
}

func int64ToInt32PtrP(i int64) *int32 {
	i32 := int32(i) // Convert int64 to int32
	return &i32     // Return pointer to int32
}
