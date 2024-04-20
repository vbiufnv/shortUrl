package url

type ByVisited []Result

// 实现 sort.Interface 接口的三个方法
func (a ByVisited) Len() int { return len(a) }

func (a ByVisited) Less(i, j int) bool { return a[i].Visited < a[j].Visited }

func (a ByVisited) Swap(i, j int) { a[i], a[j] = a[j], a[i] }