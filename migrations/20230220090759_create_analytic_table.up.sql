CREATE TABLE analytics (
    device_id   UUID,
    device_os   String,
    session     String,
    sequence    UInt64,
    event       String,
    param_int   UInt64,
    param_str   String,
    client_ip   String,
    client_time DateTime,
    server_time DateTime
) Engine=MergeTree ORDER BY device_id;

--migration:split