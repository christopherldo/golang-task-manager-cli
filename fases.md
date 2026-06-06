### **Fase 1: O Básico (Estruturas de Dados)**

Comece criando o programa para rodar apenas na memória enquanto o terminal estiver aberto.

* **O que fazer:** Crie um menu simples no terminal onde você possa adicionar uma tarefa, listar todas as tarefas e marcar uma como concluída.
* **O que você vai aprender:**
* Criação e uso de `structs` (ex: criar uma estrutura `Task` com `ID`, `Description` e `IsDone`).
* Trabalhar com *Slices* (as listas do Go) para armazenar as tarefas.
* Laços de repetição (`for`) e condicionais (`if/else`).



### **Fase 2: Persistência de Dados (Arquivos)**

Quando você fechar o programa, as tarefas vão sumir. É hora de resolver isso.

* **O que fazer:** Faça o programa salvar as tarefas em um arquivo local toda vez que houver uma alteração, e carregar esse arquivo toda vez que o programa iniciar.
* **O que você vai aprender:**
* Trabalhar com o pacote `os` para manipulação de arquivos.
* Usar o pacote `encoding/json` ou `encoding/csv` para converter suas *slices* de tarefas em texto legível.
* **Tratamento de erros:** Você vai escrever muito `if err != nil { ... }`, que é a base da segurança em Go.



### **Fase 3: Profissionalizando a CLI**

Em vez de um menu interativo, faça o programa responder a argumentos diretos no terminal.

* **O que fazer:** Fazer o programa entender comandos como:
* `meu-app add "Estudar Go"`
* `meu-app list`
* `meu-app done 1`


* **O que você vai aprender:**
* Usar o pacote nativo `flag` ou `os.Args` para capturar argumentos passados pelo usuário.
* O conceito de ponteiros (`*` e `&`), que é muito usado ao capturar *flags* no Go.



### **Fase 4 (Desafio Bônus): Transformar em uma API REST**

Depois que a sua lógica de gerenciar tarefas estiver pronta e redondinha, você pode levar o projeto para a web.

* **O que fazer:** Crie um servidor local que responda a requisições HTTP (GET, POST, PUT, DELETE) para gerenciar as tarefas, testando via Postman ou Insomnia.
* **O que você vai aprender:**
* Usar o poderoso pacote nativo `net/http`.
* Como lidar com rotas e respostas JSON na web sem precisar baixar nenhum framework externo.
