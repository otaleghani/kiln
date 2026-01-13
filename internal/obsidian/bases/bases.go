package bases

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/scanner"
	"time"

	"github.com/otaleghani/kiln/internal/obsidian"
)

// ==========================================
// 2. The API
// ==========================================

// GroupFiles organizes files into groups based on a field value.
func GroupFiles(files []*obsidian.File, field string) []*FileGroup {
	if field == "" {
		return nil
	}

	// 1. Grouping Map
	// We use a map to collect files by their group key string representation
	groupsMap := make(map[string][]*obsidian.File)

	for _, file := range files {
		val := getValue(file, field)
		key := toString(val) // Use your existing helper to convert any -> string

		// Handle empty keys (optional: label them "Uncategorized" or keep empty)
		if key == "" {
			key = "Uncategorized"
		}

		groupsMap[key] = append(groupsMap[key], file)
	}

	// 2. Convert Map to Slice
	var groups []*FileGroup
	for key, notes := range groupsMap {
		groups = append(groups, &FileGroup{
			Key:   key,
			Notes: notes,
		})
	}

	// 3. Sort Groups (Optional but recommended)
	// You might want to sort the groups themselves by Key (A-Z)
	// sort.Slice(groups, func(i, j int) bool { return groups[i].Key < groups[j].Key })

	return groups
}

// FilterFiles filters a list of files based on "and", "or", and "not" conditions.
func FilterFiles(files []*obsidian.File, filters map[string][]string) []*obsidian.File {
	// Optimization: If no filters exist, return original slice
	if len(filters) == 0 {
		return files
	}

	var filtered []*obsidian.File

	// Extract conditions once
	andConds := filters["and"]
	orConds := filters["or"]
	notConds := filters["not"]

	for _, file := range files {
		keep := true

		// 1. AND: All conditions must be true
		if len(andConds) > 0 {
			for _, query := range andConds {
				if !evalQuery(file, query) {
					keep = false
					break
				}
			}
		}
		if !keep {
			continue
		}

		// 2. NOT: None of the conditions must be true (Reject if ANY is true)
		if len(notConds) > 0 {
			for _, query := range notConds {
				if evalQuery(file, query) {
					keep = false
					break
				}
			}
		}
		if !keep {
			continue
		}

		// 3. OR: At least one condition must be true (if the block exists)
		if len(orConds) > 0 {
			orMatch := false
			for _, query := range orConds {
				if evalQuery(file, query) {
					orMatch = true
					break
				}
			}
			if !orMatch {
				keep = false
			}
		}

		if keep {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// evalQuery helps reuse the parsing logic for AND/OR/NOT loops
func evalQuery(file *obsidian.File, rawQuery string) bool {
	// 1. Parse
	p := newParser(rawQuery)
	n, err := p.parse()
	if err != nil {
		// Log error if needed, for now treat invalid syntax as false (no match)
		return false
	}

	// 2. Evaluate
	res := n.eval(file)

	// 3. Check Truthiness
	return isTrue(res)
}

// ==========================================
// 3. The AST (Abstract Syntax Tree)
// ==========================================

type node interface {
	eval(ctx *obsidian.File) any
}

type literalNode struct {
	Value any
}

func (n *literalNode) eval(ctx *obsidian.File) any { return n.Value }

type fieldNode struct {
	Name string
}

func (n *fieldNode) eval(ctx *obsidian.File) any {
	return getValue(ctx, n.Name)
}

type unaryNode struct {
	Operator string
	Operand  node
}

func (n *unaryNode) eval(ctx *obsidian.File) any {
	if n.Operand == nil {
		return false
	} // Safety check
	val := n.Operand.eval(ctx)
	switch n.Operator {
	case "!":
		return !isTrue(val)
	}
	return val
}

type methodNode struct {
	Object node
	Method string
	Args   []node // Changed to slice to support multiple args
}

func (n *methodNode) eval(ctx *obsidian.File) any {
	obj := n.Object.eval(ctx)

	// Check before
	method := strings.ToLower(n.Method)
	if method == "isempty" {
		return isEmpty(obj)
	}

	if obj == nil {
		return nil
	}

	// Evaluate all arguments
	var argValues []any
	for _, argNode := range n.Args {
		argValues = append(argValues, argNode.eval(ctx))
	}

	if method == "haslink" {
		// 1. Safety Checks
		file, ok := obj.(*obsidian.File)
		if !ok || len(argValues) == 0 {
			return false
		}

		// 2. Process the Target (User Input)
		// Input: "books/fantasy/Harry Potter.md"
		rawTarget := toString(argValues[0])

		// Remove folders: "Harry Potter.md"
		// We use LastIndex of "/" to manually strip path to ensure cross-platform safety for Obsidian paths
		if idx := strings.LastIndex(rawTarget, "/"); idx != -1 {
			rawTarget = rawTarget[idx+1:]
		}

		// Remove extension: "Harry Potter"
		target := strings.TrimSuffix(rawTarget, ".md")

		// Lowercase for fuzzy matching
		target = strings.ToLower(target)

		// 3. Search in File Links
		for _, rawLink := range file.Links {
			// Clean the link: "[[Harry Potter|HP]]" -> "Harry Potter"
			clean := cleanLink(rawLink)
			cleanLower := strings.ToLower(clean)

			// Check for containment
			// "harry potter" contains "harry potter" -> True
			// "harry potter and the stone" contains "harry potter" -> True
			if strings.Contains(cleanLower, target) {
				return true
			}
		}
		return false
	}

	// <--- ADD THIS BLOCK
	if method == "infolder" {
		// Ensure the object is actually a File
		f, ok := obj.(*obsidian.File)
		if !ok {
			return false
		}

		// Get the folder argument (e.g., "assets")
		targetFolder := ""
		if len(argValues) > 0 {
			targetFolder = toString(argValues[0])
		}

		// Logic: Check if the file's folder path contains or equals the target
		// You might want strict prefix checking or simple containment.
		// Option A: Strict folder match (file must be inside "assets" or "assets/sub")
		// We use strings.HasPrefix on the file.Folder path.
		// We add "/" to ensure we don't match "assets_backup" when looking for "assets"

		if targetFolder == "" || targetFolder == "/" {
			return true // Effectively root?
		}

		// Clean paths to be safe
		fileFolder := strings.Trim(f.Folder, "/")
		checkFolder := strings.Trim(targetFolder, "/")

		if fileFolder == checkFolder {
			return true
		}
		return strings.HasPrefix(fileFolder, checkFolder+"/")
	}
	// ---> END ADD BLOCK

	if method == "hastag" {
		// 1. Safety Check
		file, ok := obj.(*obsidian.File)
		if !ok || len(argValues) == 0 {
			return false
		}

		// 2. Prepare Target
		target := strings.ToLower(toString(argValues[0]))

		// Handle optional hash prefix: "#book" -> "book"
		target = strings.TrimPrefix(target, "#")

		// 3. Search
		for tag, _ := range file.Tags {
			// Normalize file tag: "#Book" -> "book"
			t := strings.ToLower(tag)
			t = strings.TrimPrefix(t, "#")

			// Strict Equality Check for tags is usually best
			if t == target {
				return true
			}
		}
		return false
	}

	if method == "hasproperty" {
		// 1. Safety Check
		file, ok := obj.(*obsidian.File)
		if !ok || len(argValues) == 0 {
			return false
		}

		// 2. Prepare Target
		target := strings.ToLower(toString(argValues[0]))

		// 3. Search Frontmatter Keys
		// We iterate to support case-insensitive matching (e.g. "date" finds "Date")
		for key := range file.Frontmatter {
			if strings.ToLower(key) == target {
				return true
			}
		}
		return false
	}

	// --- Multi-Argument & Collection Logic ---

	// contains / containsany / containsall
	// These can operate on a Slice Object (tags) or String Object (file.name)
	// And can accept multiple Args (A, B) or a single List Arg ([A, B])

	if strings.Contains(method, "contains") {
		// 1. Flatten Args: If first arg is a slice, use that as the target set
		var targets []any
		if len(argValues) == 1 {
			if slice, ok := toSlice(argValues[0]); ok {
				targets = slice
			} else {
				targets = []any{argValues[0]}
			}
		} else {
			targets = argValues
		}

		// 2. Determine container type
		objSlice, isObjSlice := toSlice(obj)

		switch method {
		case "contains":
			// contains(A, B) usually implies contains ANY (if multiple args provided in Obsidian?)
			// Standard Obsidian: contains(singleItem).
			// If we support multiple args here, we treat it as "contains any" or just check first?
			// Let's assume standard usage is single arg. If multiple, check if ALL?
			// To be robust: Check if object contains arg[0]
			if len(targets) == 0 {
				return false
			}

			if isObjSlice {
				// List contains Item
				return checkContains(objSlice, targets[0])
			}
			// String contains Substring
			// Normalize both the container (file.name) and the target ("harry")
			containerStr := strings.ToLower(toString(obj))
			targetStr := strings.ToLower(toString(targets[0]))

			return strings.Contains(containerStr, targetStr)
			// return strings.Contains(toString(obj), toString(targets[0]))

		case "containsany", "containsanyof": // aliases
			if isObjSlice {
				return checkAnyOf(objSlice, targets)
			}
			return checkAnyOf(toString(obj), targets)

		case "containsall", "containsallof":
			if isObjSlice {
				return checkAllOf(objSlice, targets)
			}
			return checkAllOf(toString(obj), targets)
		}
	}

	// --- Single Argument String Logic ---

	var firstArg any
	if len(argValues) > 0 {
		firstArg = argValues[0]
	}
	strObj := toString(obj)
	strArg := toString(firstArg)

	switch method {
	case "startswith":
		return strings.HasPrefix(strObj, strArg)
	case "endswith":
		return strings.HasSuffix(strObj, strArg)
	case "lower":
		return strings.ToLower(strObj)
	case "upper":
		return strings.ToUpper(strObj)
	case "length", "len":
		return len(strObj)
	}
	return obj
}

type binaryNode struct {
	Left     node
	Operator string
	Right    node
}

func (n *binaryNode) eval(ctx *obsidian.File) any {
	if n.Left == nil {
		return false
	}
	left := n.Left.eval(ctx)

	// Unary-style operators inside binaryNode
	switch n.Operator {
	case "is empty":
		return isEmpty(left)
	case "is not empty":
		return !isEmpty(left)
	}

	if n.Right == nil {
		return false
	} // Safety for binary ops
	right := n.Right.eval(ctx)

	switch n.Operator {

	case "==", "is":
		if compareValues(left, right) == 0 {
			return true
		}

		return isSameDay(left, right)
	case "!=", "is not":
		if compareValues(left, right) == 0 {
			return false
		}
		return !isSameDay(left, right)

	case ">":
		if isEmpty(left) {
			return false
		}
		return compareValues(left, right) > 0
	case ">=":
		if isEmpty(left) {
			return false
		}
		return compareValues(left, right) >= 0
	case "<":
		if isEmpty(left) {
			return false
		}
		return compareValues(left, right) < 0
	case "<=":
		if isEmpty(left) {
			return false
		}
		return compareValues(left, right) <= 0

	case "&&", "and":
		return isTrue(left) && isTrue(right)
	case "||", "or":
		return isTrue(left) || isTrue(right)

	case "on":
		if isEmpty(left) {
			return false
		}
		return isSameDay(left, right)
	case "not on":
		return !isSameDay(left, right)

	case "contains":
		// Binary Op: "tags contains 'x'" (Left contains Right)
		// Right might be a list or scalar
		return checkContainsGeneric(left, right)
	case "does not contain":
		// Logic: If left is missing/empty, should it match?
		// Standard "filtering" usually implies:
		// "Exclude files where X is present".
		// So if X is missing, it is NOT present, so we KEEP the file.
		// This is the behavior you said you disliked?

		// If you want "Must have field AND not contain":
		if isEmpty(left) {
			return false // Drop file if field is missing
		}
		return !checkContainsGeneric(left, right)

	case "contains any of":
		// Left: Container, Right: Targets
		return checkAnyOfGeneric(left, right)
	case "does not contain any of":
		if isEmpty(left) {
			return false
		}
		return !checkAnyOfGeneric(left, right)

	case "contains all of":
		return checkAllOfGeneric(left, right)
	case "does not contain all of":
		if isEmpty(left) {
			return false
		}
		return !checkAllOfGeneric(left, right)

	case "starts with":
		return strings.HasPrefix(toString(left), toString(right))
	case "does not start with":
		return !strings.HasPrefix(toString(left), toString(right))

	case "ends with":
		return strings.HasSuffix(toString(left), toString(right))
	case "does not end with":
		return !strings.HasSuffix(toString(left), toString(right))
	}
	return false
}

// ==========================================
// 4. The Parser
// ==========================================

type parser struct {
	sc   scanner.Scanner
	curr rune
	text string
}

func newParser(input string) *parser {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanInts | scanner.ScanFloats
	p := &parser{sc: s}
	p.next()
	return p
}

func (p *parser) next() {
	p.curr = p.sc.Scan()
	p.text = p.sc.TokenText()
}

func (p *parser) parse() (node, error) {
	return p.parseExpression()
}

func (p *parser) parseExpression() (node, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for {
		op := strings.ToLower(p.text)
		if op == "and" || op == "&&" || op == "or" || op == "||" {
			p.next()
			right, err := p.parseComparison()
			if err != nil {
				return nil, err
			}
			left = &binaryNode{Left: left, Operator: op, Right: right}
		} else {
			break
		}
	}
	return left, nil
}

func (p *parser) parseComparison() (node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	op := strings.ToLower(p.text)
	var right node
	unary := false

	// Handle Operators
	if op == "is" {
		p.next()
		if strings.ToLower(p.text) == "not" {
			p.next()
			if strings.ToLower(p.text) == "empty" {
				op = "is not empty"
				unary = true
				p.next()
			} else {
				op = "is not"
			}
		} else if strings.ToLower(p.text) == "empty" {
			op = "is empty"
			unary = true
			p.next()
		}

	} else if op == "does" {
		p.next()
		if strings.ToLower(p.text) == "not" {
			p.next()
			sub := strings.ToLower(p.text)
			if sub == "contain" {
				// ... (Parsing logic for "any of", "all of" etc)
				op = "does not contain"
			} else if sub == "start" {
				p.next()
				if strings.ToLower(p.text) == "with" {
					p.next()
					op = "does not start with"
				}
			} else if sub == "end" {
				p.next()
				if strings.ToLower(p.text) == "with" {
					p.next()
					op = "does not end with"
				}
			} else if sub == "contain" {
				p.next()
				if strings.ToLower(p.text) == "any" {
					p.next()
					if strings.ToLower(p.text) == "of" {
						p.next()
						op = "does not contain any of"
					}
				} else if strings.ToLower(p.text) == "all" {
					p.next()
					if strings.ToLower(p.text) == "of" {
						p.next()
						op = "does not contain all of"
					}
				} else {
					op = "does not contain"
				}
			}
		}

	} else if op == "contains" {
		p.next()
		if strings.ToLower(p.text) == "any" {
			p.next()
			if strings.ToLower(p.text) == "of" {
				p.next()
				op = "contains any of"
			}
		} else if strings.ToLower(p.text) == "all" {
			p.next()
			if strings.ToLower(p.text) == "of" {
				p.next()
				op = "contains all of"
			}
		}

	} else if op == "starts" {
		p.next()
		if strings.ToLower(p.text) == "with" {
			p.next()
			op = "starts with"
		}
	} else if op == "ends" {
		p.next()
		if strings.ToLower(p.text) == "with" {
			p.next()
			op = "ends with"
		}

	} else if isOperator(p.text) {
		op = p.text
		p.next()
		if p.text == "=" {
			op += "="
			p.next()
		}
	} else {
		return left, nil
	}

	if !unary {
		right, err = p.parseUnary()
		if err != nil {
			return nil, err
		}
	}

	return &binaryNode{Left: left, Operator: op, Right: right}, nil
}

func (p *parser) parseUnary() (node, error) {
	if p.text == "!" {
		p.next()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &unaryNode{Operator: "!", Operand: operand}, nil
	}
	return p.parseTerm()
}

func (p *parser) parseTerm() (node, error) {
	n, err := p.parseFactor() // Renamed variable 'node' to 'n' to avoid shadowing type 'node'
	if err != nil {
		return nil, err
	}

	for p.text == "." {
		p.next()

		part := p.text
		p.next()

		if p.text == "(" {
			p.next()

			// Parse Arguments (Multi-arg support)
			var args []node
			for p.text != ")" && p.curr != scanner.EOF {
				arg, err := p.parseExpression()
				if err != nil {
					return nil, err
				} // Check error to prevent nil nodes
				args = append(args, arg)

				if p.text == "," {
					p.next()
				}
			}

			if p.text == ")" {
				p.next()
			}

			n = &methodNode{Object: n, Method: part, Args: args}
		} else {
			if fNode, ok := n.(*fieldNode); ok {
				fNode.Name = fNode.Name + "." + part
				n = fNode
			} else {
				// Fallback for complex property access if needed
				return nil, fmt.Errorf("unexpected property access")
			}
		}
	}
	return n, nil
}

func (p *parser) parseFactor() (node, error) {
	tok := p.curr
	text := p.text

	switch tok {
	case scanner.Ident:
		if text == "true" {
			p.next()
			return &literalNode{Value: true}, nil
		}
		if text == "false" {
			p.next()
			return &literalNode{Value: false}, nil
		}
		if text == "null" || text == "nil" {
			p.next()
			return &literalNode{Value: nil}, nil
		}
		p.next()
		return &fieldNode{Name: text}, nil

	case scanner.String:
		val, _ := strconv.Unquote(text)
		p.next()
		return &literalNode{Value: val}, nil

	case scanner.Int, scanner.Float:
		f, _ := strconv.ParseFloat(text, 64)
		p.next()
		return &literalNode{Value: f}, nil

	case '-':
		p.next()
		if p.curr == scanner.Int || p.curr == scanner.Float {
			f, _ := strconv.ParseFloat(p.text, 64)
			p.next()
			return &literalNode{Value: -f}, nil
		}

	case '[':
		p.next()
		var list []any
		for p.text != "]" && p.curr != scanner.EOF {
			item, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			if lit, ok := item.(*literalNode); ok {
				list = append(list, lit.Value)
			}

			if p.text == "," {
				p.next()
			}
		}
		p.next()
		return &literalNode{Value: list}, nil

	case '(':
		p.next()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.text == ")" {
			p.next()
		}
		return expr, nil
	}
	return nil, fmt.Errorf("unexpected token: %s", text)
}

func isOperator(s string) bool {
	return s == "=" || s == "!" || s == ">" || s == "<" || s == "==" || s == "!="
}

// ==========================================
// 5. Helpers
// ==========================================

func getValue(file *obsidian.File, field string) any {
	switch field {
	// --- Identity ---
	case "file":
		return file
	case "file.name":
		return file.Name // e.g. "MyNote.md"
	case "file.stem", "file.basename":
		// Just the name without extension (e.g. "MyNote")
		name := file.Name
		if idx := strings.LastIndex(name, "."); idx != -1 {
			return name[:idx]
		}
		return name
	case "file.path":
		return file.Path // The file system path
	case "file.link":
		return file.WebPath // The URL/Web path

	// --- Dates ---
	case "file.ctime":
		return file.Created // Returns time.Time object
	case "file.mtime":
		return file.Modified // Returns time.Time object
	case "file.cday":
		// New: Useful for grouping by Date (ignoring time)
		return file.Created.Format("2006-01-02")
	case "file.mday":
		return file.Modified.Format("2006-01-02")

	// --- Metadata ---
	case "file.size":
		return file.Size
	case "file.ext":
		return strings.TrimPrefix(file.Ext, ".")

	// --- Collections (The missing parts) ---
	case "file.links":
		return file.Links // []string
	case "file.embeds":
		return file.Embeds // []string
	case "file.tags":
		// Vital: Convert map[string]struct{} to []string
		// Otherwise filters like 'contains(file.tags, ...)' will fail.
		var tags []string
		for t := range file.Tags {
			tags = append(tags, t)
			// Obsidian accepts both with or without "#" for tags
			tags = append(tags, strings.TrimPrefix(t, "#"))
		}
		return tags
	case "file.folder":
		return file.Folder
	}

	// Fallback: Check Frontmatter (YAML)

	if val, ok := file.Frontmatter[field]; ok {
		return val
	}
	return nil
}

func isSameDay(a, b any) bool {
	t1, ok1 := toTime(a)
	t2, ok2 := toTime(b)

	// If both are valid times/dates
	if ok1 && ok2 {
		y1, m1, d1 := t1.Date()
		y2, m2, d2 := t2.Date()
		return y1 == y2 && m1 == m2 && d1 == d2
	}

	// Fallback: If they aren't dates, try string comparison
	// (matches logic if user does: file.name on "MyFile")
	return toString(a) == toString(b)
}

// compareValues compares two values and returns:
// -1 if a < b
//
//	0 if a == b
//	1 if a > b
func compareValues(a, b any) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// 1. Try Numbers (Int/Float)
	f1, ok1 := toFloat(a)
	f2, ok2 := toFloat(b)
	if ok1 && ok2 {
		if f1 > f2 {
			return 1
		}
		if f1 < f2 {
			return -1
		}
		return 0
	}

	// 2. Try Dates (Time)
	// This handles: file.ctime > "2023-01-01"
	t1, isTime1 := toTime(a)
	t2, isTime2 := toTime(b)

	// We only compare as time if AT LEAST one of them was originally a true time.Time object
	// or if both successfully parsed as dates.
	// This prevents random strings like "10-5" from being treated as dates unless intended.
	if isTime1 && isTime2 {
		if t1.After(t2) {
			return 1
		}
		if t1.Before(t2) {
			return -1
		}
		return 0
	}

	// 3. Fallback to Strings
	s1 := toString(a)
	s2 := toString(b)
	if s1 > s2 {
		return 1
	}
	if s1 < s2 {
		return -1
	}
	return 0
}

// toTime attempts to convert any value to a time.Time
func toTime(v any) (time.Time, bool) {
	// If it's already a Time object (from file.ctime or YAML date)
	if t, ok := v.(time.Time); ok {
		return t, true
	}

	// If it's a string, try to parse it
	if s, ok := v.(string); ok {
		// List of formats to try.
		// You can add more formats here if needed (e.g. European "02-01-2006")
		formats := []string{
			"2006-01-02",                // YYYY-MM-DD (Common Obsidian/YAML)
			"2006-01-02 15:04:05",       // SQL / Standard
			"2006-01-02T15:04:05Z07:00", // RFC3339 (ISO)
			"2006-01-02 15:04",          // Short datetime
		}

		for _, f := range formats {
			if t, err := time.Parse(f, s); err == nil {
				return t, true
			}
		}
	}

	return time.Time{}, false
}

func isEmpty(v any) bool {
	// 1. Nil is always empty
	if v == nil {
		return true
	}

	// 2. Strings
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s) == ""
	}

	// 3. Reflection for Slices, Maps, Arrays
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		return val.Len() == 0
	}

	// 4. Fallback (Numbers/Bools are never "empty" in this context)
	return false
}

// Wrapper for Binary Ops to handle types safely
// func checkContainsGeneric(container, item any) bool {
// 	if slice, ok := toSlice(container); ok {
// 		return checkContains(slice, item)
// 	}
// 	return strings.Contains(toString(container), toString(item))
// }

func checkContainsGeneric(container, item any) bool {
	// 1. If it's a List (Links, Tags, Embeds) -> Use Slice Logic
	if slice, ok := toSlice(container); ok {
		return checkContains(slice, item)
	}

	// 2. If it's a String (Name, Content) -> Use String Logic (Case-Insensitive)
	containerStr := strings.ToLower(toString(container))
	itemStr := strings.ToLower(toString(item))
	return strings.Contains(containerStr, itemStr)
}

func checkAnyOfGeneric(left, right any) bool {
	// 1. Ensure 'right' is a slice. If not, wrap it.
	targets, ok := toSlice(right)
	if !ok {
		targets = []any{right}
	}

	// 2. If 'left' is a slice (e.g. file.tags), check against list items
	if slice, ok := toSlice(left); ok {
		return checkAnyOf(slice, targets)
	}

	// 3. If 'left' is a string (e.g. file.name), check substrings
	return checkAnyOf(toString(left), targets)
}

func checkAllOfGeneric(left, right any) bool {
	// 1. Ensure 'right' is a slice
	targets, ok := toSlice(right)
	if !ok {
		targets = []any{right}
	}

	if slice, ok := toSlice(left); ok {
		return checkAllOf(slice, targets)
	}
	return checkAllOf(toString(left), targets)
}

// Low-level Logic

// checkContains handles: file.tags.contains("x") OR file.links.contains("x")
// func checkContains(slice []any, item any) bool {
// 	// Lowercase the target once
// 	target := strings.ToLower(toString(item))
//
// 	for _, v := range slice {
// 		rawS := toString(v)
// 		cleanS := cleanLink(rawS)
//
// 		// Check 1: Exact Raw Match (Case-Insensitive)
// 		// Matches "[[Link]]" against "[[link]]"
// 		if strings.ToLower(rawS) == target {
// 			return true
// 		}
//
// 		// Check 2: Cleaned Strict Match (Case-Insensitive)
// 		// Matches "[[Harry Potter|HP]]" against "harry potter"
// 		if strings.ToLower(cleanS) == target {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// // checkAnyOf handles: file.tags.containsAny("x", "y")
// func checkAnyOf(left any, targets []any) bool {
// 	// Case A: Left is a Slice (links, tags, embeds)
// 	if source, ok := left.([]any); ok {
// 		for _, s := range source {
// 			rawS := toString(s)
// 			cleanS := cleanLink(rawS)
//
// 			lowerRaw := strings.ToLower(rawS)
// 			lowerClean := strings.ToLower(cleanS)
//
// 			for _, t := range targets {
// 				target := strings.ToLower(toString(t))
//
// 				if lowerRaw == target {
// 					return true
// 				}
// 				if lowerClean == target {
// 					return true
// 				}
// 			}
// 		}
// 		return false
// 	}
//
// 	// Case B: Left is a String (file.name, file.content)
// 	strLeft := strings.ToLower(toString(left))
// 	for _, t := range targets {
// 		if strings.Contains(strLeft, strings.ToLower(toString(t))) {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// // checkAllOf handles: file.tags.containsAll("x", "y")
// func checkAllOf(left any, targets []any) bool {
// 	// Case A: Left is a Slice
// 	if source, ok := left.([]any); ok {
// 		for _, t := range targets {
// 			found := false
// 			target := strings.ToLower(toString(t))
//
// 			for _, s := range source {
// 				rawS := toString(s)
// 				cleanS := cleanLink(rawS)
//
// 				if strings.ToLower(rawS) == target || strings.ToLower(cleanS) == target {
// 					found = true
// 					break
// 				}
// 			}
// 			if !found {
// 				return false
// 			}
// 		}
// 		return true
// 	}
//
// 	// Case B: Left is a String
// 	strLeft := strings.ToLower(toString(left))
// 	for _, t := range targets {
// 		if !strings.Contains(strLeft, strings.ToLower(toString(t))) {
// 			return false
// 		}
// 	}
// 	return true
// }

// checkContains: Case-Sensitive for Slices
func checkContains(slice []any, item any) bool {
	// Target is kept AS IS (Case-Sensitive)
	target := toString(item)

	for _, v := range slice {
		rawS := toString(v)
		cleanS := cleanLink(rawS)

		// Check 1: Exact Raw Match (Case-Sensitive)
		// "[[Page]]" == "[[Page]]" (True)
		// "[[Page]]" == "[[page]]" (False)
		if rawS == target {
			return true
		}

		// Check 2: Cleaned Strict Match (Case-Sensitive)
		// "[[Page]]" cleans to "Page".
		// "Page" == "Page" (True)
		// "Page" == "page" (False)
		if cleanS == target {
			return true
		}
	}
	return false
}

// checkAnyOf: Strict for Slices, Fuzzy for Strings
func checkAnyOf(left any, targets []any) bool {
	// Case A: Left is a Slice (Links/Embeds) -> Strict / Case-Sensitive
	if source, ok := left.([]any); ok {
		for _, s := range source {
			rawS := toString(s)
			cleanS := cleanLink(rawS)

			for _, t := range targets {
				target := toString(t) // No ToLower

				if rawS == target {
					return true
				}
				if cleanS == target {
					return true
				}
			}
		}
		return false
	}

	// Case B: Left is a String (Name/Text) -> Fuzzy / Case-Insensitive
	strLeft := strings.ToLower(toString(left))
	for _, t := range targets {
		if strings.Contains(strLeft, strings.ToLower(toString(t))) {
			return true
		}
	}
	return false
}

// checkAllOf: Strict for Slices, Fuzzy for Strings
func checkAllOf(left any, targets []any) bool {
	// Case A: Left is a Slice -> Strict / Case-Sensitive
	if source, ok := left.([]any); ok {
		for _, t := range targets {
			found := false
			target := toString(t) // No ToLower

			for _, s := range source {
				rawS := toString(s)
				cleanS := cleanLink(rawS)

				if rawS == target || cleanS == target {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}

	// Case B: Left is a String -> Fuzzy / Case-Insensitive
	strLeft := strings.ToLower(toString(left))
	for _, t := range targets {
		if !strings.Contains(strLeft, strings.ToLower(toString(t))) {
			return false
		}
	}
	return true
}

// cleanLink normalizes Obsidian links/embeds for comparison.
// Input: "![[My Page#Section|Alias]]" -> Output: "My Page"
func cleanLink(s string) string {
	// 1. Remove Embed bang "!"
	s = strings.TrimPrefix(s, "!")

	// 2. Remove Brackets "[[" and "]]"
	s = strings.TrimPrefix(s, "[[")
	s = strings.TrimSuffix(s, "]]")

	// 3. Remove Alias (stop at "|")
	if idx := strings.Index(s, "|"); idx != -1 {
		s = s[:idx]
	}

	// 4. Remove Anchor/Heading (stop at "#")
	if idx := strings.Index(s, "#"); idx != -1 {
		s = s[:idx]
	}

	return strings.TrimSpace(s)
}

func isTrue(v any) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

//	func toString(v any) string {
//		if v == nil {
//			return ""
//		}
//		return fmt.Sprintf("%v", v)
//	}

// toString converts arbitrary values into a consistent string key.
func toString(val any) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		// Format floats to remove trailing zeros if needed
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case time.Time:
		return v.Format("2006-01-02") // Standardize dates for grouping

	// Handle the Tag map specifically
	case map[string]struct{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys) // Sort to ensure "a,b" is the same group as "b,a"
		return strings.Join(keys, ", ")

	// Handle standard slices (if Frontmatter contains lists)
	case []any:
		strs := make([]string, len(v))
		for i, item := range v {
			strs[i] = toString(item)
		}
		return strings.Join(strs, ", ")

	default:
		return fmt.Sprintf("%v", v)
	}
}

func toFloat(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float64:
		return val, true
	}
	return 0, false
}

func toSlice(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice {
		out := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			out[i] = val.Index(i).Interface()
		}
		return out, true
	}
	return nil, false
}
