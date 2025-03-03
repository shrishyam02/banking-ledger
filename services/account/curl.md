curl -X POST -H "Content-Type: application/json" -u test:test -d '{
  "accountNumber": "1234567890982",
  "balance": 1000.00,
  "customer": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john.dab@example.com",
    "phoneNumber": "1234567890982",
    "address": "123 Main St"
  }
}' http://localhost:8000/api/v1/accounts


curl -X POST -H "Content-Type: application/json" -u test:test -d '{
  "accountNumber": "12345678900",
  "balance": 1000.00,
  "accountType": "savings"
}' http://localhost:8001/api/v1/accounts

curl -u test:test http://localhost:8001/api/v1/accounts/8db6626d-5e84-4c4e-8cec-7dc54cb20ff5


curl -X POST -H "Content-Type: application/json" -u test:test -d '{
  "accountID": "6e032122-ef4a-4cc6-a531-61b8a9ee7320",
  "amount": 1000.00,
  "transactionType": "credit"
}' http://localhost:7002/api/v1/transactions









$ curl -X POST -H "Content-Type: application/json" -u test:test -d '{
  "accountNumber": "1234567890989",
  "balance": 1000.00,
  "customerId": "00000000-0000-0000-0000-000000000000",
  "customer": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john.da@example.com",
    "phoneNumber": "1234567890989",
    "address": "123 Main St"
  }
}' http://localhost:8000/api/v1/accounts
{"ID":"8db6626d-5e84-4c4e-8cec-7dc54cb20ff5","AccountNumber":"1234567890989","AccountType":"","Status":"active","Balance":1000,"CreatedAt":"2025-03-01T17:47:09.768777525Z","UpdatedAt":"2025-03-01T17:47:09.768777525Z","CustomerID":"550e8400-e29b-41d4-a716-446655440000","Customer":{"ID":"550e8400-e29b-41d4-a716-446655440000","Name":"John Doe","Email":"john.da@example.com","PhoneNumber":"1234567890989","Address":"123 Main St","CreatedAt":"2025-03-01T17:47:09.766031884Z","UpdatedAt":"2025-03-01T17:47:09.749901916Z"}}


curl -X POST -H "Content-Type: application/json" -u test:test -d '{
  "accountID": "8db6626d-5e84-4c4e-8cec-7dc54cb20ff5",
  "amount": 1000.00,
  "transactionType": "credit"
}' http://localhost:7002/api/v1/transactions



curl -X POST -H "Content-Type: application/json" -u test:test -d '{
  "accountID": "8db6626d-5e84-4c4e-8cec-7dc54cb20ff5",
  "amount": 1000.00,
  "transactionType": "credit"
}' http://localhost:8000/api/v1/transactions



curl -X POST -H "Content-Type: application/json" -u test:test -d '{
  "accountID": "8db6626d-5e84-4c4e-8cec-7dc54cb20ff5",
  "amount": 1000.00,
  "transactionType": "credit"
}' http://localhost:7002/api/v1/transactions


curl -u test:test http://localhost:7004/api/v1/ledger/accounts/8db6626d-5e84-4c4e-8cec-7dc54cb20ff5



curl -u test:test http://localhost:7004/api/v1/ledger/accounts/8db6626d-5e84-4c4e-8cec-7dc54cb20ff5
[{"_id":"67c55f3f4d1761dd15c0d70d","acceptedAt":"2025-03-03T07:50:22.996Z","accountId":"8db6626d-5e84-4c4e-8cec-7dc54cb20ff5","amount":50,"details":"debit","id":"152a42be-63b6-46f9-919e-ba3996eaa890","processedAt":"2025-03-03T07:50:22.996Z","status":"success","transactionType":"debit"},{"_id":"67c55f3f4d1761dd15c0d70c","acceptedAt":"2025-03-03T07:50:22.996Z","accountId":"8db6626d-5e84-4c4e-8cec-7dc54cb20ff5","amount":100,"details":"Initial deposit","id":"d6263bc8-0eeb-4195-9e64-81abd6d5685c","processedAt":"2025-03-03T07:50:22.996Z","status":"success","transactionType":"credit"}]

curl -u test:test http://localhost:7004/api/v1/ledger/transactions/152a42be-63b6-46f9-919e-ba3996eaa890
[{"_id":"67c55f3f4d1761dd15c0d70d","acceptedAt":"2025-03-03T07:50:22.996Z","accountId":"8db6626d-5e84-4c4e-8cec-7dc54cb20ff5","amount":50,"details":"debit","id":"152a42be-63b6-46f9-919e-ba3996eaa890","processedAt":"2025-03-03T07:50:22.996Z","status":"success","transactionType":"debit"}]