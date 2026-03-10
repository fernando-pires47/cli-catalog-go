## Design - <nome-da-feature>

### 1) Objetivo

<objetivo tecnico do desenho>

### 2) Principios de Arquitetura

- <principio 1>
- <principio 2>

### 3) Visao de Camadas

- `app/*`: <...>
- `components/ui/*`: <...>
- `components/<dominio>/*`: <...>
- `hooks/*`: <...>
- `services/*`: <...>
- `types/*`: <...>
- `utils/*`: <...>

### 4) Rotas e Responsabilidades

- `/<rota-1>`: <...>
- `/<rota-2>`: <...>

### 5) Modelagem de Estado (Frontend)

#### 5.1 Estado Global

- Campos: <...>
- Acoes: <...>

#### 5.2 Estado Local

- <formularios/loading/erro>

### 6) Contratos de Dados (Types)

```ts
interface Example {
  id: string
}
```

### 7) Design de Servicos

- `<service>`: `<metodo>(input): Promise<output>`

### 8) Regras de Negocio no Design

- <regra implementada>

### 9) Estrategia de Assincronia

- Polling/retry/cache: <...>
- Paradas de fluxo: <...>

### 10) UX/UI e Theming

- Tokens centralizados
- Estados: loading, empty, erro, sucesso

### 11) Validacao e Formularios

- <biblioteca e schemas>

### 12) Observabilidade e Telemetria

- <evento 1>
- <evento 2>

### 13) Seguranca e Privacidade

- <diretriz 1>

### 14) Estrategia de Testes

- Unitarios: <...>
- Integracao: <...>
- E2E: <...>

### 15) Plano de Entrega Incremental

- Fase 1: <...>
- Fase 2: <...>

### 16) Riscos Tecnicos e Mitigacoes

- <risco> -> <mitigacao>

### 17) Decisoes em Aberto

- <decisao pendente>
