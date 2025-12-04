package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Input struct {
	Params []struct {
		InputName string `json:"inputname"`
		CompValue string `json:"compvalue"`
	} `json:"params"`
}

type Output struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

func main() {
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		json.NewEncoder(os.Stdout).Encode(Output{Error: fmt.Sprintf("failed to decode input: %v", err)})
		return
	}

	var (
		decisionType string
		leftValue    string
		rightValue   string
		operator     string
	)

	for _, p := range input.Params {
		switch p.InputName {
		case "decisiontype":
			decisionType = strings.ToLower(strings.TrimSpace(p.CompValue))
		case "leftvalue":
			leftValue = strings.TrimSpace(p.CompValue)
		case "rightvalue":
			rightValue = strings.TrimSpace(p.CompValue)
		case "operator":
			operator = strings.TrimSpace(p.CompValue)
		}
	}

	var decision bool
	var err error

	switch decisionType {
	case "numeric":
		decision, err = evaluateNumeric(leftValue, rightValue, operator)
	case "string":
		decision = evaluateString(leftValue, rightValue, operator)
	case "boolean":
		decision, err = evaluateBoolean(leftValue)
	default:
		json.NewEncoder(os.Stdout).Encode(Output{Error: fmt.Sprintf("unknown decision type: %s", decisionType)})
		return
	}

	if err != nil {
		json.NewEncoder(os.Stdout).Encode(Output{Error: err.Error()})
		return
	}

	result := map[string]interface{}{
		"decision": decision,
		"message":  fmt.Sprintf("%s %s %s = %v", leftValue, operator, rightValue, decision),
	}

	json.NewEncoder(os.Stdout).Encode(Output{Result: result})
}

func evaluateNumeric(left, right, operator string) (bool, error) {
	leftNum, err := strconv.ParseFloat(left, 64)
	if err != nil {
		return false, fmt.Errorf("invalid left numeric value: %s", left)
	}

	rightNum, err := strconv.ParseFloat(right, 64)
	if err != nil {
		return false, fmt.Errorf("invalid right numeric value: %s", right)
	}

	switch operator {
	case "==", "=":
		return leftNum == rightNum, nil
	case "!=", "<>":
		return leftNum != rightNum, nil
	case ">":
		return leftNum > rightNum, nil
	case ">=":
		return leftNum >= rightNum, nil
	case "<":
		return leftNum < rightNum, nil
	case "<=":
		return leftNum <= rightNum, nil
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}

func evaluateString(left, right, operator string) bool {
	switch operator {
	case "==", "=":
		return left == right
	case "!=", "<>":
		return left != right
	case "contains":
		return strings.Contains(left, right)
	case "startswith":
		return strings.HasPrefix(left, right)
	case "endswith":
		return strings.HasSuffix(left, right)
	default:
		return false
	}
}

func evaluateBoolean(value string) (bool, error) {
	value = strings.ToLower(value)
	switch value {
	case "true", "1", "yes", "y":
		return true, nil
	case "false", "0", "no", "n", "":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", value)
	}
}
