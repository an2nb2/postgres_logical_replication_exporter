INSERT INTO samples1 (uid, data)
SELECT LEFT(md5(i::text), 10), md5(random()::text)
FROM generate_series(1, 1000) s(i);

INSERT INTO samples2 (uid, data)
SELECT LEFT(md5(i::text), 10), md5(random()::text)
FROM generate_series(1, 1000) s(i);
