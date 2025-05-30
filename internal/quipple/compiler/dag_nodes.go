package compiler

import (
	"qp/internal/pkgdata"
	"sync"
)

type FilterNode interface {
	Stream(input <-chan *pkgdata.PkgInfo) <-chan *pkgdata.PkgInfo
}

type QueryNode struct {
	Filter func(*pkgdata.PkgInfo) bool
}

func (n *QueryNode) Stream(input <-chan *pkgdata.PkgInfo) <-chan *pkgdata.PkgInfo {
	output := make(chan *pkgdata.PkgInfo)

	go func() {
		defer close(output)

		for pkg := range input {
			if n.Filter(pkg) {
				output <- pkg
			}
		}
	}()

	return output
}

type AndNode struct {
	Left  FilterNode
	Right FilterNode
}

func (n *AndNode) Stream(input <-chan *pkgdata.PkgInfo) <-chan *pkgdata.PkgInfo {
	return n.Right.Stream(n.Left.Stream(input))
}

type OrNode struct {
	Left  FilterNode
	Right FilterNode
}

func (n *OrNode) Stream(input <-chan *pkgdata.PkgInfo) <-chan *pkgdata.PkgInfo {
	output := make(chan *pkgdata.PkgInfo)

	go func() {
		defer close(output)

		leftIn := make(chan *pkgdata.PkgInfo)
		rightIn := make(chan *pkgdata.PkgInfo)

		go func() {
			for pkg := range input {
				leftIn <- pkg
				rightIn <- pkg
			}

			close(leftIn)
			close(rightIn)
		}()

		seen := make(map[string]bool)
		var mu sync.Mutex

		merge := func(in <-chan *pkgdata.PkgInfo) {
			for pkg := range in {
				mu.Lock()
				if !seen[pkg.Name] {
					seen[pkg.Name] = true
					output <- pkg
				}

				mu.Unlock()
			}
		}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			merge(n.Left.Stream(leftIn))
		}()

		go func() {
			defer wg.Done()
			merge(n.Right.Stream(rightIn))
		}()

		wg.Wait()
	}()

	return output
}

type NotNode struct {
	Inner FilterNode
}

func (n *NotNode) Stream(input <-chan *pkgdata.PkgInfo) <-chan *pkgdata.PkgInfo {
	output := make(chan *pkgdata.PkgInfo)

	go func() {
		defer close(output)

		var buffer []*pkgdata.PkgInfo
		inputCopy := make(chan *pkgdata.PkgInfo)

		go func() {
			defer close(inputCopy)

			for pkg := range input {
				buffer = append(buffer, pkg)
				inputCopy <- pkg
			}
		}()

		accepted := make(map[string]bool)
		for pkg := range n.Inner.Stream(inputCopy) {
			accepted[pkg.Name] = true
		}

		for _, pkg := range buffer {
			if !accepted[pkg.Name] {
				output <- pkg
			}
		}
	}()

	return output
}
