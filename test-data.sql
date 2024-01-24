INSERT INTO public.configurations (id, created_at, updated_at, deleted_at, version) VALUES (1, '2024-01-22 20:33:13.319000 +00:00', '2024-01-22 20:33:13.319000 +00:00', null, 1);
INSERT INTO public.configurations (id, created_at, updated_at, deleted_at, version) VALUES (2, '2024-01-22 20:33:42.194000 +00:00', '2024-01-22 20:33:42.194000 +00:00', null, 2);

INSERT INTO public.parameters (id, created_at, updated_at, deleted_at, configuration_id, name, type, comparer) VALUES (1, '2024-01-22 20:34:02.314000 +00:00', '2024-01-22 20:34:02.314000 +00:00', null, 1, 'color', 'string', 'eq');
INSERT INTO public.parameters (id, created_at, updated_at, deleted_at, configuration_id, name, type, comparer) VALUES (2, '2024-01-22 20:35:22.978000 +00:00', '2024-01-22 20:35:22.978000 +00:00', null, 1, 'brand', 'string', 'eq');
INSERT INTO public.parameters (id, created_at, updated_at, deleted_at, configuration_id, name, type, comparer) VALUES (3, '2024-01-22 20:36:02.879000 +00:00', '2024-01-22 20:36:02.879000 +00:00', null, 1, 'horsepower', 'int', 'gt');
INSERT INTO public.parameters (id, created_at, updated_at, deleted_at, configuration_id, name, type, comparer) VALUES (4, '2024-01-22 20:36:14.709000 +00:00', '2024-01-22 20:36:14.709000 +00:00', null, 1, 'year', 'datetime', 'gt');
INSERT INTO public.parameters (id, created_at, updated_at, deleted_at, configuration_id, name, type, comparer) VALUES (5, '2024-01-22 20:34:02.314000 +00:00', '2024-01-22 20:34:02.314000 +00:00', null, 2, 'color', 'string', 'eq');
INSERT INTO public.parameters (id, created_at, updated_at, deleted_at, configuration_id, name, type, comparer) VALUES (6, '2024-01-22 20:35:22.978000 +00:00', '2024-01-22 20:35:22.978000 +00:00', null, 2, 'brand', 'string', 'eq');
INSERT INTO public.parameters (id, created_at, updated_at, deleted_at, configuration_id, name, type, comparer) VALUES (7, '2024-01-22 20:36:02.879000 +00:00', '2024-01-22 20:36:02.879000 +00:00', null, 2, 'horsepower', 'int', 'gt');
INSERT INTO public.parameters (id, created_at, updated_at, deleted_at, configuration_id, name, type, comparer) VALUES (8, '2024-01-22 20:36:14.709000 +00:00', '2024-01-22 20:36:14.709000 +00:00', null, 2, 'year', 'datetime', 'gt');
INSERT INTO public.parameters (id, created_at, updated_at, deleted_at, configuration_id, name, type, comparer) VALUES (9, '2024-01-22 20:36:14.709000 +00:00', '2024-01-22 20:36:14.709000 +00:00', null, 2, 'consumption', 'float', 'le');

INSERT INTO public.products (id, created_at, updated_at, deleted_at, configuration_id, name) VALUES (1, '2024-01-24 18:37:01.853000 +00:00', '2024-01-24 18:37:01.853000 +00:00', null, 1,'SlowCar');
INSERT INTO public.products (id, created_at, updated_at, deleted_at, configuration_id, name) VALUES (2, '2024-01-24 18:37:01.853000 +00:00', '2024-01-24 18:37:01.853000 +00:00', null, 1,'FastCar');
INSERT INTO public.products (id, created_at, updated_at, deleted_at, configuration_id, name) VALUES (3, '2024-01-24 18:37:01.853000 +00:00', '2024-01-24 18:37:01.853000 +00:00', null, 2,'SlowCar');
INSERT INTO public.products (id, created_at, updated_at, deleted_at, configuration_id, name) VALUES (4, '2024-01-24 18:37:01.853000 +00:00', '2024-01-24 18:37:01.853000 +00:00', null, 2,'FastCar');

INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (1, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 1, 1,'Red');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (2, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 2, 1,'Ferrari');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (3, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 3, 1,'150');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (4, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 4, 1,'1970-01-01 00:00:00.000000 +00:00');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (5, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 1, 2,'Red');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (6, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 2, 2,'Ferrari');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (7, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 3, 2,'250');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (8, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 4, 2,'1970-01-01 00:00:00.000000 +00:00');


INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (9, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 5, 3,'Red');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (10, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 6, 3,'Ferrari');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (11, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 7, 3,'150');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (12, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 8, 3,'1970-01-01 00:00:00.000000 +00:00');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (13, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 9, 3,'7.8');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (14, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 5, 4,'Red');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (15, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 6, 4,'Ferrari');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (16, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 7, 4,'250');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (17, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 8, 4,'1970-01-01 00:00:00.000000 +00:00');
INSERT INTO public.parameter_values (id, created_at, updated_at, deleted_at, parameter_id, product_id, value) VALUES (18, '2024-01-24 18:39:45.520000 +00:00', '2024-01-24 18:39:45.520000 +00:00', null, 9, 4,'11.3');

