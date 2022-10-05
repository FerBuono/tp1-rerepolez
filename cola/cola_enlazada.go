package cola

type nodoCola[T any] struct {
	dato T
	prox *nodoCola[T]
}

type colaEnlazada[T any] struct {
	primero *nodoCola[T]
	ultimo  *nodoCola[T]
}

func (c *colaEnlazada[T]) EstaVacia() bool {
	return c.primero == nil && c.ultimo == nil
}

func (c *colaEnlazada[T]) VerPrimero() T {
	if c.EstaVacia() {
		panic("La cola esta vacia")
	}

	return c.primero.dato
}

func (c *colaEnlazada[T]) Encolar(dato T) {
	nuevo := c.crearNodo(dato)

	if c.EstaVacia() {
		c.primero = nuevo
	} else {
		c.ultimo.prox = nuevo
	}

	c.ultimo = nuevo
}

func (c *colaEnlazada[T]) Desencolar() T {
	if c.EstaVacia() {
		panic("La cola esta vacia")
	}

	elemento := c.primero.dato

	if c.primero.prox == nil {
		c.ultimo = nil
	}
	c.primero = c.primero.prox

	return elemento
}

func (c *colaEnlazada[T]) crearNodo(dato T) *nodoCola[T] {
	nuevo := new(nodoCola[T])
	nuevo.dato = dato
	return nuevo
}

func CrearColaEnlazada[T any]() Cola[T] {
	c := new(colaEnlazada[T])
	return c
}
