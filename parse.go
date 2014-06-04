package html

import (
	"fmt"
	"html/template"
	"text/template/parse"
)

const (
	rootTemplateName  = ""
	dummyTemplateName = "__dummy__"
)

func createTemplate(parsedTrees []*parse.Tree, funcs map[string]interface{}) (*template.Template, error) {
	tmpl, _ := template.New(dummyTemplateName).Funcs(funcs).Parse("")
	for _, tree := range parsedTrees {
		var err error
		tmpl, err = tmpl.AddParseTree(tree.Name, tree)
		if err != nil {
			return nil, err
		}
	}

	err := validateTemplate(tmpl)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func validateTemplate(root *template.Template) error {
	if root.Lookup(rootTemplateName) == nil {
		return fmt.Errorf("html/template: missing root template")
	}

	needTemplates := []string{}
	haveTemplates := make(map[string]bool)
	for _, tmpl := range root.Templates() {
		tmplName := tmpl.Tree.Name
		if tmplName != "" {
			haveTemplates[tmplName] = true
		}
		needTemplates = append(needTemplates, requiredTemplates(tmpl.Tree.Root)...)
	}

	var missingTemplates []string
	for _, tmplName := range needTemplates {
		if found := haveTemplates[tmplName]; !found {
			missingTemplates = append(missingTemplates, tmplName)
		}
	}

	if len(missingTemplates) > 0 {
		return fmt.Errorf("html/template: missing template(s) %q", missingTemplates)
	}

	return nil
}

func requiredTemplates(root *parse.ListNode) (names []string) {
	if root == nil {
		return
	}

	for _, node := range root.Nodes {
		if tnode, ok := node.(*parse.TemplateNode); ok {
			names = append(names, tnode.Name)
		} else if lnode, ok := node.(*parse.ListNode); ok {
			names = append(names, requiredTemplates(lnode)...)
		} else if bnode, ok := node.(*parse.IfNode); ok {
			names = append(names, requiredTemplates(bnode.BranchNode.List)...)
			names = append(names, requiredTemplates(bnode.BranchNode.ElseList)...)
		} else if bnode, ok := node.(*parse.RangeNode); ok {
			names = append(names, requiredTemplates(bnode.BranchNode.List)...)
			names = append(names, requiredTemplates(bnode.BranchNode.ElseList)...)
		}
	}
	return
}
