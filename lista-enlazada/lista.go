package lista

type Lista[T any] interface {
	// EstaVacia devuelve true si la lista no tiene elementos, false en caso contrario.
	EstaVacia() bool

	// InsertarPrimero agrega un nuevo elemento en el primer lugar de la lista.
	InsertarPrimero(T)

	// InsertarUltimo agrega un nuevo elemento en el último lugar de la lista.
	InsertarUltimo(T)

	// BorrarPrimero elimina el primer elemento de la lista. Si la lista tiene elementos, se quita el primero de la misma,
	// y se devuelve su valor. Si está vacía, entra en pánico con un mensaje "La lista esta vacia".
	BorrarPrimero() T

	// VerPrimero devuelve el valor del primer elemento de la lista. Si la lista tiene elementos, se muestra el valor
	// del primero. Si está vacía, entra en pánico con un mensaje "La lista esta vacia".
	VerPrimero() T

	// VerUltimo devuelve el valor del último elemento de la lista. Si la lista tiene elementos, se muestra el valor
	// del último. Si está vacía, entra en pánico con un mensaje "La lista esta vacia".
	VerUltimo() T

	// Largo devuelve la cantidad de elementos que hay en la lista. Si está vacía devuelve 0.
	Largo() int

	// Iterador devuelve un iterador de la lista, el cual cuenta con sus primitivas.
	Iterador() IteradorLista[T]

	// Iterar recibe una funcion que se va a aplicar los datos de la lista de manera ordenada, hasta que se acabe
	// la lista o hasta que dicha funcion devuelva false
	Iterar(func(T) bool)
}

type IteradorLista[T any] interface {
	// VerActual devuelve el valor del elemento en el que se encuentra el iterador. En caso de que el iterador ya haya
	// iterado todos los elementos, entra en pánico con un mensaje "El iterador termino de iterar".
	VerActual() T

	// HaySiguiente devuelve true si existe un elemento siguiente al actual, false en caso contrario.
	HaySiguiente() bool

	// Siguiente devuelve el valor del elemento en el que se encuentra el iterador, y luego avanza al siguiente.
	// En caso de que el iterador ya haya iterado todos los elementos, entra en pánico con un mensaje "El iterador termino de iterar".
	Siguiente() T

	// Insertar añade un nuevo elemento a la lista en la posición en la que se encuntra el iterador, moviendo el elemento
	// que se encontraba originalmente en esa posición a la siguiente.
	Insertar(T)

	// Borrar elimina el elemento que se encuentra en la posición actual del iterador, vinculando el anterior con el siguiente.
	// En caso de que el iterador ya haya iterado todos los elementos, entra en pánico con un mensaje "El iterador termino de iterar".
	Borrar() T
}
