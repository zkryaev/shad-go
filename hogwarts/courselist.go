//go:build !solution

package hogwarts

import (
	"fmt"
	//"slices"
)

const (
	unvisited = iota
	visited
	parented
)

func DFS(graph map[string][]string, v string, used *map[string]int, ans *[]string) {
	stack := make([]string, 0, len(*used))
	(*used)[v] = visited
	cycle := make([]string, 0, len(*used))
	stack = append(stack, v)
	cycle = append(cycle, v)
	for len(stack) > 0 {
		// DFS function
		from := stack[len(stack)-1]
		if (*used)[from] == parented || len(graph[from]) == 0 {
			stack = stack[:len(stack)-1]
			*ans = append(*ans, from)
			(*used)[from] = visited
		} else {
			(*used)[from] = parented
			for _, to := range graph[from] {
				if (*used)[to] == parented {
					cycle = append(cycle, to)
					panic(cycle)
				}
				if (*used)[to] == unvisited {
					(*used)[to] = visited
					stack = append(stack, to)
					cycle = append(cycle, to)
				}
			}
		}
	}
}

func topologicSort(graph map[string][]string) []string {
	ans := make([]string, 0, len(graph))
	used := make(map[string]int)
	for v := range graph {
		used[v] = unvisited
	}
	for v := range graph {
		if used[v] == unvisited {
			DFS(graph, v, &used, &ans)
		}
	}
	fmt.Println("result:", ans)
	//slices.Reverse(ans)
	return ans
}

func GetCourseList(prereqs map[string][]string) []string {
	// check cycle
	// topologicSort
	CourseList := topologicSort(prereqs)

	return CourseList
}
