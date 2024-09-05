package concurrency

import "sync"

// ExecuteInParallel executa uma função (action) em paralelo com controle de limites e etapas.
//
// A função é projetada para dividir a execução de tarefas em "etapas" (steps), onde cada etapa
// pode executar múltiplas instâncias da função `action` em paralelo até atingir o limite especificado.
// A execução de cada instância é gerenciada por um `sync.WaitGroup` para garantir que todas as tarefas
// em uma etapa sejam concluídas antes que a próxima etapa comece.
//
// Parâmetros:
// - limit: O número total de vezes que a função `action` será executada. Este valor define o limite máximo de execuções.
// - steps: O número de execuções paralelas em cada etapa. Controla quantas execuções simultâneas da função `action` ocorrerão por vez.
// - action: A função que será executada em paralelo. Recebe como argumento um índice (idx) que varia de 0 até `limit-1`.
//
// Comportamento:
// - A função executa `action` em paralelo em blocos definidos por `steps`.
// - Em cada etapa, até `steps` execuções de `action` são iniciadas em goroutines paralelas.
// - A função garante que todas as execuções em uma etapa sejam concluídas antes de iniciar a próxima etapa.
//
// Exemplo de uso:
// ```go
//
//	ExecuteInParallel(10, 3, func(idx int) {
//	    fmt.Printf("Executando ação %d\n", idx)
//	})
//
// ```
// Neste exemplo, a função `action` será executada 10 vezes, em lotes de 3 execuções paralelas por vez.
// Isso significa que em cada etapa, até 3 execuções ocorrerão simultaneamente até que todas as 10 execuções sejam concluídas.
func ExecuteInParallel(limit, steps int, action func(idx int)) {
	var wg sync.WaitGroup
	for i := 0; i < limit; i += steps {
		for j := 0; j < steps; j++ {

			idx := i + j
			if idx >= limit {
				break
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				action(idx)
			}()
		}
		wg.Wait()
	}
}
