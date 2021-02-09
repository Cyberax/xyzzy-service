-- name: Ping :one
select 1;

-- name: InsertDataset :execrows
insert into test(id, name, url) values (@id, @name, @url);

-- name: SelectDatasetByName :many
select * from test where (name = @name);
