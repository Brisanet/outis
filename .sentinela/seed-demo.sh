#!/bin/bash

# --- CONFIGURAÇÕES ---
OTEL_URL="http://localhost:4318/v1"
LOOPS=20 # Quantidade de transações a gerar
ERROR_RATE=5 # 1 a cada X requisições será um erro

# Função para gerar HEX (Linux nativo)
generate_hex() {
    cat /dev/urandom | tr -dc 'a-f0-9' | fold -w $1 | head -n 1
}

echo "🚀 Iniciando simulação de tráfego intenso..."
echo "📦 Cenário: Checkout de E-commerce (Gateway -> Auth -> Payment -> DB)"
echo "🔄 Gerando $LOOPS transações..."

for ((i=1; i<=LOOPS; i++)); do
    # 1. Gerar IDs únicos para esta transação
    TRACE_ID=$(generate_hex 32)
    SPAN_GW=$(generate_hex 16)
    SPAN_AUTH=$(generate_hex 16)
    SPAN_PAY=$(generate_hex 16)
    SPAN_DB=$(generate_hex 16)
    
    # 2. Definir Timestamps (Cascata temporal)
    NOW_NANO=$(date +%s%N)
    T_GW_START=$NOW_NANO
    T_AUTH_START=$(($NOW_NANO + 50000000))   # +50ms
    T_AUTH_END=$(($NOW_NANO + 150000000))     # +150ms
    T_PAY_START=$(($NOW_NANO + 200000000))    # +200ms
    T_DB_START=$(($NOW_NANO + 250000000))     # +250ms
    T_DB_END=$(($NOW_NANO + 450000000))       # +450ms
    T_PAY_END=$(($NOW_NANO + 500000000))      # +500ms
    T_GW_END=$(($NOW_NANO + 550000000))       # +550ms

    # 3. Roleta Russa de Erros (Simular falha no pagamento)
    IS_ERROR=0
    STATUS_CODE=1 # 1 = OK no OTLP
    LOG_LEVEL="INFO"
    LOG_MSG="Transação aprovada com sucesso"
    
    if (( i % ERROR_RATE == 0 )); then
        IS_ERROR=1
        STATUS_CODE=2 # 2 = ERROR no OTLP
        LOG_LEVEL="ERROR"
        LOG_MSG="Falha ao processar cartão de crédito: Saldo insuficiente"
        echo "   ⚠️  Gerando erro simulado na transação $i..."
    fi

    # ---------------------------------------------------------
    # ENVIO DE TRACES (Batch com 4 serviços interligados)
    # ---------------------------------------------------------
    curl -s -X POST "$OTEL_URL/traces" -H "Content-Type: application/json" -d '{
     "resourceSpans": [
       {
         "resource": { "attributes": [{ "key": "service.name", "value": { "stringValue": "api-gateway" } }, { "key": "service_name", "value": { "stringValue": "api-gateway" } }] },
         "scopeSpans": [{
           "spans": [{
             "traceId": "'$TRACE_ID'",
             "spanId": "'$SPAN_GW'",
             "name": "POST /checkout",
             "kind": 1, 
             "startTimeUnixNano": "'$T_GW_START'",
             "endTimeUnixNano": "'$T_GW_END'",
             "status": { "code": 1 }
           }]
         }]
       },
       {
         "resource": { "attributes": [{ "key": "service.name", "value": { "stringValue": "auth-service" } }, { "key": "service_name", "value": { "stringValue": "auth-service" } }] },
         "scopeSpans": [{
           "spans": [{
             "traceId": "'$TRACE_ID'",
             "spanId": "'$SPAN_AUTH'",
             "parentSpanId": "'$SPAN_GW'",
             "name": "Validate Token",
             "kind": 2,
             "startTimeUnixNano": "'$T_AUTH_START'",
             "endTimeUnixNano": "'$T_AUTH_END'",
             "status": { "code": 1 }
           }]
         }]
       },
       {
         "resource": { "attributes": [{ "key": "service.name", "value": { "stringValue": "payment-service" } }, { "key": "service_name", "value": { "stringValue": "payment-service" } }] },
         "scopeSpans": [{
           "spans": [{
             "traceId": "'$TRACE_ID'",
             "spanId": "'$SPAN_PAY'",
             "parentSpanId": "'$SPAN_GW'",
             "name": "Process Credit Card",
             "kind": 2,
             "startTimeUnixNano": "'$T_PAY_START'",
             "endTimeUnixNano": "'$T_PAY_END'",
             "status": { "code": '$STATUS_CODE', "message": "'"$LOG_MSG"'" }
           }]
         }]
       },
       {
         "resource": { "attributes": [{ "key": "service.name", "value": { "stringValue": "postgres-db" } }, { "key": "service_name", "value": { "stringValue": "postgres-db" } }] },
         "scopeSpans": [{
           "spans": [{
             "traceId": "'$TRACE_ID'",
             "spanId": "'$SPAN_DB'",
             "parentSpanId": "'$SPAN_PAY'",
             "name": "INSERT INTO orders",
             "kind": 3,
             "startTimeUnixNano": "'$T_DB_START'",
             "endTimeUnixNano": "'$T_DB_END'",
             "attributes": [{ "key": "db.statement", "value": { "stringValue": "INSERT INTO orders VALUES (...)" } }]
           }]
         }]
       }
     ]
    }' > /dev/null

    # ---------------------------------------------------------
    # ENVIO DE LOGS (Correlacionados via TraceID e ServiceName)
    # ---------------------------------------------------------
    # Log do Gateway (Sempre Info)
    curl -s -X POST "$OTEL_URL/logs" -H "Content-Type: application/json" -d '{
      "resourceLogs": [{
        "resource": { "attributes": [] }, 
        "scopeLogs": [{
          "logRecords": [{
            "timeUnixNano": "'$T_GW_START'",
            "severityText": "INFO",
            "body": { "stringValue": "{\"message\": \"Iniciando checkout de pedido\", \"service_name\": \"api-gateway\", \"trace_id\": \"'$TRACE_ID'\", \"http_method\": \"POST\"}" }
          }]
        }]
      }]
    }' > /dev/null

    # Log do Payment (Pode ser Erro ou Info)
    # Note que enviamos o JSON stringificado dentro do body para seu Transform processar
    curl -s -X POST "$OTEL_URL/logs" -H "Content-Type: application/json" -d '{
      "resourceLogs": [{
        "resource": { "attributes": [] }, 
        "scopeLogs": [{
          "logRecords": [{
            "timeUnixNano": "'$T_PAY_END'",
            "severityText": "'$LOG_LEVEL'",
            "body": { "stringValue": "{\"message\": \"'$LOG_MSG'\", \"service_name\": \"payment-service\", \"trace_id\": \"'$TRACE_ID'\", \"amount\": \"150.00\"}" }
          }]
        }]
      }]
    }' > /dev/null

    echo -ne "   ✅ Transação $i/$LOOPS enviada (TraceID: ${TRACE_ID:0:8}...)\r"
    sleep 0.2 # Pequena pausa para não floodar instantaneamente
done

echo -e "\n\n🎉 Simulação concluída!"
echo "👉 Vá ao Grafana e verifique:"
echo "   1. O dropdown 'Serviço' deve ter 4 novas opções."
echo "   2. O Service Graph (Explore -> Tempo) deve mostrar a árvore de dependências."
echo "   3. Busque logs de 'payment-service' e procure pelas falhas."

# Fonte: Gemini PRO