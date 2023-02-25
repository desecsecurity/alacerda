package main

func main() {
	_, rSub, _ := matematica(1, 4)
	println(rSub)
}

func matematica(num1, num2 int) (soma, subtracao, multiplicacao int) {
	soma = num1 + num2
	subtracao = num1 - num2
	multiplicacao = num1 * num2

	return
}
