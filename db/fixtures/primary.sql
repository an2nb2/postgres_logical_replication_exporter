CREATE ROLE exporter WITH LOGIN PASSWORD '1234';

CREATE PUBLICATION pub1 FOR TABLE public.samples1;
CREATE PUBLICATION pub2 FOR TABLE public.samples2;
