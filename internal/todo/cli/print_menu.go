package cli

import "fmt"

func printMenuCli(session ProgramSession) {
	switch session {
	case MenuSession:
		fmt.Println(`===============================================
O que deseja?
1 - Adicionar uma task.
2 - Listar todas as tasks adicionadas.
3 - Marcar uma task como concluída.
4 - Editar uma task.
5 - Deletar uma task.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
	case AddedTaskSession:
		fmt.Println(`===============================================
O que deseja?
1 - Adicionar outra task.
2 - Listar todas as tasks adicionadas.
3 - Marcar uma task como concluída.
4 - Editar uma task.
5 - Deletar uma task.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
	case ListAllTasksSession:
		fmt.Println(`===============================================
O que deseja agora?
1 - Adicionar uma task.
2 - Listar novamente todas as tasks adicionadas.
3 - Marcar uma task como concluída.
4 - Editar uma task.
5 - Deletar uma task.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
	case CompletingTaskSession:
		fmt.Println(`===============================================
O que deseja agora?
1 - Adicionar uma task.
2 - Listar todas as tasks adicionadas.
3 - Marcar outra task como concluída.
4 - Editar uma task.
5 - Deletar uma task.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
	case EditingTaskSession:
		fmt.Println(`===============================================
O que deseja agora?
1 - Adicionar uma task.
2 - Listar todas as tasks adicionadas.
3 - Marcar uma task como concluída.
4 - Editar outra task.
5 - Deletar uma task.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
	case DeletingTaskSession:
		fmt.Println(`===============================================
O que deseja agora?
1 - Adicionar uma task.
2 - Listar todas as tasks adicionadas.
3 - Marcar uma task como concluída.
4 - Editar outra task.
5 - Deletar outra task.
9 - Voltar ao menu principal.
0 - Sair do programa.`)
	default:
		break
	}
}
