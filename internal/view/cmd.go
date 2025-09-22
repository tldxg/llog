package view

import "fmt"

func PrintGet(count int) {
	fmt.Println(entryAction(count, "retrieved"))
}

func PrintAdd(count int) {
	fmt.Println(entryAction(count, "created"))
}

func PrintDelete(count int) {
	fmt.Println(entryAction(count, "deleted"))
}

func PrintSummarize(count int) {
	fmt.Println(entryAction(count, "summarized"))
}

func entryAction(count int, action string) string {
	switch count {
	case 0:
		return fmt.Sprintf("no entry %s", action)
	case 1:
		return fmt.Sprintf("1 entry %s", action)
	default:
		return fmt.Sprintf("%d entries %s", count, action)
	}
}

func PrintNuke() {
	fmt.Println("database nuked")
}
