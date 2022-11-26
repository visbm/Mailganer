# avito-tech
1. Start app in docker 
`docker-compose up --build`

2. Import postman collection `avito-tech.postman_collection.json`

3. Requests in Postman :
    3.1 GET localhost:8080/get/all 
            description: Return all accounts
                response:
                    body:
                        application/json:
                        [
                                {
                                    "account_id": {"type":"int"},
                                    "balance": {"type":"int"}
                                },
                                {
                                    "account_id": {"type":"int"},
                                    "balance": {"type":"int"}
                                },
                            ]
                        example:
                            [
                                {
                                    "account_id": 1,
                                    "balance": 333
                                },
                                {
                                    "account_id":2,
                                    "balance": 33.5
                                },
                            ]

    3.2 GET localhost:8080/get/balance/{id} 
            description: Return account by id
                example: localhost:8080/get/balance/1
                    response:
                        body:
                            application/json:                        
                                {
                                    "account_id": {"type":"int"},
                                    "balance": {"type":"int"}
                                }
                                example:
                                {
                                    "account_id": 2,
                                    "balance": 444
                                }
                            
    3.3 POST localhost:8080/transaction 
             description: make transaction from one user to ther
                request:
                  body:
                    application/json:
                        {
                            "receiver_id":{"type":"int"},   //required
                            "sender_id": {"type":"int"},    //required
                            "amount": {"type":"int"},       //required
                            "comment": {"type":"string"}    //required
                        }
                    example:
                        {
                            "receiver_id": 2,
                            "sender_id": 1,
                            "amount": 20,
                            "comment": "You paid for me in a restaurant"
                        }
                response: 
                    TEXT:
                        "Transaction was successful"

    3.4 POST localhost:8080/changeBalance                       
             description: chane ballance vy account id
                request:    
                  body:
                    application/json:
                        {
                            "account_id":{"type":"int"},    //required
                            "amount": {"type":"int"},       //required
                            "comment": {"type":"string"}    //required
                        }
                    example:
                        {
                            "account_id": 2,
                            "amount": 20,
                            "comment": "Salary"
                        }   
                response: 
                    TEXT:
                        "Balance succsefully changed"

    3.5 GET localhost:8080/get/balance/{currency}/{id} 
            description: Return account balance in currency by id                    
                example: localhost:8080/get/balance/USD/1
                response:
                    body:
                            application/json:                        
                                {
                                   "result": {"type":"int"}
                                }
                                example:
                                {
                                    "result": 444
                                }

    3.6 GET localhost:8080/get/transactions            
            description: chane ballance vy account id
                request:
                  body:
                    application/json:
                        {
                            "account_id":{"type":"int"},    //required
                            "limit":{"type":"int"},         
                            "offset": {"type":"int"}, 
                            "order_by" : {"type":"string"}  //IN : ("date_time", "amount","date_time ASC", "amount ASC","date_time DESC", "amount DESC","date_time asc", "amount asc","date_time desc", "amount desc")
                        }
                    example:
                        { 
                            "account_id": 1, 
                            "limit":0,
                            "offset": 0, 
                            "order_by" : "date_time DESC"  
                        }   
                response:
                        body:
                            application/json:                        
                                [
                                        {
                                            "transaction_ID": {"type":"int"},
                                            "account_id": {"type":"int"},
                                            "amount": {"type":"int"},
                                            "date": {"type":"string"},
                                            "comment": {"type":"string"}
                                        },
                                        {
                                            "transaction_ID": {"type":"int"},
                                            "account_id": {"type":"int"},
                                            "amount": {"type":"int"},
                                            "date": {"type":"string"},
                                            "comment": {"type":"string"}
                                        }
                                ]        
                                example:
                                [
                                        {
                                            "transaction_ID": 1,
                                            "account_id": 1,
                                            "amount": -20,
                                            "date": "2022-10-04T08:12:00.508216Z",
                                            "comment": "You paid for me in a restaurant"
                                        },
                                        {
                                            "transaction_ID": 2,
                                            "account_id": 1,
                                            "amount": 165.5,
                                            "date": "1999-02-08T00:00:00Z",
                                            "comment": "Salary"
                                        }
                                ]        