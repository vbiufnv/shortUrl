package url

type ByVisited []Result

// 实现 sort.Interface 接口的 Len 方法
func (a ByVisited) Len() int { return len(a) }

// 实现 sort.Interface 接口的 Less 方法
func (a ByVisited) Less(i, j int) bool { return a[i].Visited < a[j].Visited }

// 实现 sort.Interface 接口的 Swap 方法
func (a ByVisited) Swap(i, j int) { a[i], a[j] = a[j], a[i] }