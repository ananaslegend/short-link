CREATE TABLE IF NOT EXISTS statistic (
    redirect_time_stamp DateTime,
    link text,
    redirect UInt64
)
ENGINE = MergeTree()
ORDER BY redirect_time_stamp;