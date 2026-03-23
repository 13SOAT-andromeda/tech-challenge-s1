# Fluxo de Criação e Ciclo de Vida de uma Order

## Pré-requisitos

Antes de criar uma order, os seguintes recursos precisam existir no sistema:

```mermaid
flowchart TD
    subgraph PRE["Pré-requisitos"]
        A([Criar Usuário\nPOST /users]) --> B([Criar Empresa\nPOST /companies])
        C([Criar Cliente\nPOST /customers]) --> D([Criar Veículo\nPOST /vehicles])
        D --> E([Associar Veículo ao Cliente\nPOST /customers/:id/vehicles/:vehicleId])
        F([Criar Produto\nPOST /products])
        G([Criar Manutenção\nPOST /maintenances])
    end

    PRE --> H([Criar Order\nPOST /orders])
```

---

## Criação da Order

```mermaid
sequenceDiagram
    actor Admin as Administrador
    participant Router as Router<br/>(AuthRequired + RoleRequired)
    participant Handler as OrderHandler
    participant UseCase as OrderUseCase
    participant Service as OrderService
    participant DB as PostgreSQL

    Admin->>Router: POST /orders<br/>Authorization: Bearer JWT
    Router->>Router: Valida JWT (AuthRequired)
    Router->>Router: Valida role=administrator (RoleRequired)
    Router->>Handler: CreateOrderRequest<br/>{ vehicle_kilometers, note?,<br/>customer_vehicle_id, company_id }

    Handler->>Handler: ShouldBindJSON (valida campos obrigatórios)
    Handler->>Handler: Extrai user_id do contexto JWT

    Handler->>UseCase: CreateOrder(ctx, userID, input)

    UseCase->>UseCase: Monta domain.Order<br/>Status = "Recebida"<br/>DateIn = agora()

    UseCase->>Service: Create(ctx, order)
    Service->>Service: model.FromDomain(order)
    Service->>DB: INSERT INTO orders
    DB-->>Service: order criada (com ID)
    Service->>Service: model.ToDomain()
    Service-->>UseCase: *domain.Order
    UseCase-->>Handler: *domain.Order
    Handler-->>Admin: 201 Created<br/>{ success, data: Order, message }
```

---

## Ciclo de Vida Completo da Order

```mermaid
stateDiagram-v2
    [*] --> Recebida : POST /orders\n(Create)

    Recebida --> EmDiagnóstico : POST /orders/:id/assign\n(Assign — atribui técnico)

    EmDiagnóstico --> DiagnósticoFinalizado : POST /orders/:id/complete-analysis\n(CompleteAnalysis — associa produtos,\nmanutencoes e calcula preço)

    DiagnósticoFinalizado --> AguardandoAprovação : POST /orders/:id/request-approval\n(RequestApproval — envia e-mail ao cliente)

    AguardandoAprovação --> Aprovado : GET /orders/:id/approve\n(cliente aprova via link no e-mail)
    AguardandoAprovação --> Finalizado : GET /orders/:id/reject\n(cliente rejeita via link no e-mail)

    Aprovado --> EmExecução : POST /orders/:id/start-work\n(StartWork — decrementa estoque)

    EmExecução --> Finalizado : POST /orders/:id/complete-work\n(CompleteWork)

    Finalizado --> Entregue : POST /orders/:id/archive\n(Archive — registra DateOut)

    Entregue --> [*]
```

---

## Fluxo Detalhado Passo a Passo

```mermaid
flowchart TD
    START([Início]) --> P1

    subgraph STEP1["1. Criar Order"]
        P1["POST /orders\n{ vehicle_kilometers, customer_vehicle_id, company_id, note? }\nStatus → Recebida"]
    end

    subgraph STEP2["2. Atribuir Técnico"]
        P2["POST /orders/:id/assign\n{ user_id }\nStatus → Em diagnóstico"]
    end

    subgraph STEP3["3. Concluir Diagnóstico"]
        P3["POST /orders/:id/complete-analysis\n{ diagnostic_note?, products: [{id, qty}], maintenances: [id] }\nCalcula preço total\nStatus → Diagnóstico finalizado"]
    end

    subgraph STEP4["4. Solicitar Aprovação"]
        P4["POST /orders/:id/request-approval\nEnvia e-mail ao cliente com link de aprovação\nStatus → Aguardando aprovação"]
    end

    subgraph STEP5["5. Resposta do Cliente"]
        P5A["GET /orders/:id/approve\n(link público no e-mail)\nStatus → Aprovado\nregistra DateApproved"]
        P5B["GET /orders/:id/reject\n(link público no e-mail)\nStatus → Finalizado\nregistra DateRejected"]
    end

    subgraph STEP6["6. Iniciar Trabalho"]
        P6["POST /orders/:id/start-work\nDecrementa estoque dos produtos\nStatus → Em execução"]
    end

    subgraph STEP7["7. Concluir Trabalho"]
        P7["POST /orders/:id/complete-work\nStatus → Finalizado"]
    end

    subgraph STEP8["8. Arquivar"]
        P8["POST /orders/:id/archive\nRegistra DateOut\nStatus → Entregue"]
    end

    STEP1 --> STEP2
    STEP2 --> STEP3
    STEP3 --> STEP4
    STEP4 --> STEP5
    P5A --> STEP6
    P5B --> REJECTED([Encerrado\nrejeitado])
    STEP6 --> STEP7
    STEP7 --> STEP8
    STEP8 --> END([Fim])
```

---

## Estrutura do Payload de Criação

```mermaid
classDiagram
    class CreateOrderRequest {
        +int vehicle_kilometers
        +string note (opcional)
        +uint customer_vehicle_id
        +uint company_id
    }

    class Order {
        +uint id
        +time date_in
        +time date_out (nullable)
        +time date_approved (nullable)
        +time date_rejected (nullable)
        +string status
        +int vehicle_kilometers
        +string note (nullable)
        +string diagnostic_note (nullable)
        +float64 price (nullable)
        +uint customer_vehicle_id
        +uint user_id
        +uint company_id
    }

    class CustomerVehicle {
        +uint id
        +uint customer_id
        +uint vehicle_id
    }

    class Vehicle {
        +uint id
        +string plate
        +string name
        +int year
        +string brand
        +string color
    }

    class Customer {
        +uint id
        +string name
        +string email
        +string document
        +string contact
    }

    class Company {
        +uint id
        +string name
    }

    CreateOrderRequest --> Order : cria
    Order --> CustomerVehicle : referencia
    CustomerVehicle --> Vehicle : contém
    CustomerVehicle --> Customer : pertence a
    Order --> Company : pertence a
```

---

## Regras de Negócio

| Etapa | Regra |
|---|---|
| Criação | Status inicial sempre `Recebida`. `DateIn` = momento da criação. |
| Atribuição | Apenas ordens com status `Recebida` podem ser atribuídas. |
| Diagnóstico | Preço calculado automaticamente: `Σ(produto.price × qty) + Σ(manutencao.price)`. |
| Aprovação | E-mail enviado ao cliente com links públicos (sem JWT) para aprovar ou rejeitar. |
| Início do trabalho | Estoque dos produtos é decrementado **somente** nesta etapa. |
| Arquivamento | Registra `DateOut`. Ordens `Finalizado` e `Entregue` são excluídas do `GET /orders`. |
