package db


func CreatAllTable()error{
    err:=CreateSubscriptionsTable()
    if err!=nil{
        return err
    }
    err=CreatWebHookTable()
    if err!=nil{
        return err
    }
    err=CreateLogTable()
    if err!=nil{
        return err
    }
    return nil
}




func CreateSubscriptionsTable()error{
	query:=`CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    target_url TEXT NOT NULL,
    secret TEXT,
    event_types TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
   _,err:=DB.Exec(query)
   return err
}

func CreatWebHookTable()error{
    query:=`CREATE TABLE IF NOT EXISTS webhooks (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,  -- Using JSONB for flexible and variable payload
    subscription_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',  
    delivered boolean default false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    retries INT DEFAULT 0 ,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
);`
    _,err:=DB.Exec(query)
    return err
}

func CreateLogTable()error{
	query:=`CREATE TABLE IF NOT EXISTS delivery_logs (
    id SERIAL PRIMARY KEY,
    webhook_id UUID NOT NULL,
    subscription_id UUID NOT NULL,
    target_url VARCHAR(255) NOT NULL,
    attempt_number INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    http_status_code INT,
    error_details TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (webhook_id) REFERENCES webhooks(id) ON DELETE CASCADE,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
   );`  
   _,err:=DB.Exec(query)
   return err
}
