# Desafio Técnico: Processamento Concorrente de Mensagens

## Descrição do Desafio

Você receberá um banco de dados SQLite contendo uma tabela chamada `messages_table` com os seguintes campos:

- **ID**: Identificador único da mensagem.
- **data**: Campo JSON armazenado utilizando o recurso JSON1 do SQLite. Este campo contém as seguintes informações:
    - `worker`: Identificador do worker (1 a 5).
    - `message`: Texto da mensagem.
    - `interval`: Intervalo em milisegundos até a próxima mensagem.
    - `destination_worker`: Identificador do worker de destino.

### Consulte a documentação do JSON1 do SQLite:
- [JSON1 - SQLite](https://www.sqlite.org/json1.html)

### Exemplo da Estrutura da Tabela

| ID  | data   |
|-----|--------|
| 1   | ...    |
| 2   | ...    |
| 3   | ...    |
| 4   | ...    |

### Exemplo do Conteúdo do Campo `data`

```json
{
    "worker": 1,
    "message": "Olá, como vai?",
    "interval": 3,
    "destination_worker": 2
}
```

## Sua Tarefa

- **Leitura e Processamento**: Ler o banco de dados fornecido e processar as mensagens de acordo com as informações contidas no campo `data`.
- **Exibição de Mensagens**: Exibir as mensagens em qualquer formato de saída que preferir, respeitando os intervalos de tempo (`interval`) que definem a sequência de mensagens.
- **Informações a Serem Impressas**:
    - `worker`: Identificador do worker que enviou a mensagem.
    - `worker_listened`: Identificador do último worker que falou com ele (ou 0 se for o primeiro).
    - `message`: A mensagem enviada.
    - `message_listened`: A mensagem recebida do worker anterior (ou vazio se for o primeiro).

### Exemplo de Saída Esperada

```
worker: 1
worker_listened: 0
message: "Olá, como vai?"
message_listened: ""

worker: 2
worker_listened: 1
message: "Tudo bem, e você? Pode falar com o worker 5?"
message_listened: "Olá, como vai?"

worker: 1
worker_listened: 2
message: "Olá, como vai? A worker 1 pediu para falar contigo."
message_listened: "Tudo bem, e você? Pode falar com o worker 5?"

worker: 1
worker_listened: 2
message: "De novo! Vou bem, obrigada!"
message_listened: "Tudo bem, e você? Pode falar com o worker 5?"

...
```

## Orientações

- **Número de Workers**: Existem apenas 5 workers, identificados pelos números 1 a 5.
- **Linguagem de Programação**: Você é livre para resolver o desafio na linguagem de programação de sua preferência.
    - **Dica**: Uma solução elegante em Go (Golang) poderia envolver "Raw Queries" (Raw & Scan) a criação de 5 goroutines, uma para cada worker, utilizando channels para comunicação.
- **Objetivo**: Demonstrar sua habilidade em lidar com concorrência e comunicação entre processos ou threads, além de sua capacidade de ler e interpretar dados JSON armazenados em um banco de dados SQLite.

## O que Será Fornecido

- **Banco de Dados**: Um arquivo SQLite já populado com os dados necessários.
- **Script de Criação**: O código com Go e SQL que foi utilizado para criar e popular a tabela `messages_table`.
    - **Dica**: O código pode te ajudar a solucionar o desafio, mas não é necessário utilizá-lo.

# Atenção
## Utilize o Banco de Dados SQLite Fornecido
##### O script não irá recriar o banco exatamente como está no arquivo fornecido. Iremos avaliar o output baseado no Banco de Dados fornecido.


## Instruções de Entrega

- **Código-Fonte**: Envie o código-fonte completo de sua solução.
- **Instruções**: Inclua instruções claras sobre como executar seu programa.
- **Documentação**: Comente seu código quando necessário para explicar sua lógica e decisões tomadas.

## Dicas Adicionais

- **Qualidade do Código**: Preocupe-se com a legibilidade e a organização do seu código.
- **Tratamento de Erros**: Implemente um tratamento adequado de erros e exceções.
- **Testes**: Se possível, inclua testes que demonstrem o funcionamento correto do seu programa.

## Segunda Fase do Teste

- **Antecipação**: Haverá uma segunda fase do teste relacionada a noções de UI e UX, que será realizada em texto.
- **Sugestão**: Ao desenvolver sua solução, considere mentalmente como seria uma interface para exibir esses resultados. Não é necessário implementar nada relacionado a interface agora.

## Prazo

- **Tempo Estimado**: Espera-se que o desafio seja concluído em até 3 horas.
- **Flexibilidade**: Sabemos que pode ser difícil encaixar o tempo para resolver o desafio em sua rotina atual por isso você terá até o [prazo comunicado por e-mail] para nos enviar a solução.

## Como Enviar Sua Solução

Por favor, envie seu código e as instruções de execução para o seguinte e-mail:

- **E-mail**: adm@healthgo.com.br
- **Assunto**: Desafio Técnico - [Seu Nome]

---

Desejamos boa sorte!

---

HealthGo© All rights reserved. Do not copy, distribute or modify without permission. 2024

Todos os direitos reservados à HealthGo Technologies. Não copie, distribua ou modifique sem permissão. 2024