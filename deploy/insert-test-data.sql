INSERT INTO lamps (title, power_w, luminous_flux_lm, scattering_angle_deg, image_url, is_deleted) VALUES
  ('Умная потолочная люстра с перламутром', 40, 3200, 120, 'http://localhost:9000/lamp-images/1.jpg', false),
  ('Slim Magnetic Трековый светильник 26W 4000K Most чёрный', 26, 2210, 60, 'http://localhost:9000/lamp-images/2.jpg', false),
  ('Светильник потолочный светодиодный Trio 8W 3000K белый', 8, 680, 120, 'http://localhost:9000/lamp-images/3.jpg', false),
  ('Потолочный светильник', 40, 3200, 120, 'http://localhost:9000/lamp-images/4.jpg', false),
  ('Подвесной светильник со стеклянными плафонами', 40, 3200, 120, 'http://localhost:9000/lamp-images/5.jpg', false),
  ('Esthetic Magnetic Трековый светильник 3W 3000K (чёрный)', 3, 255, 60, 'http://localhost:9000/lamp-images/6.jpg', false),
  ('Подвесной светильник', 40, 3200, 120, 'http://localhost:9000/lamp-images/7.jpg', false),
  ('Трековый светильник 100W 4200K Full Light N05 Slim Magnetic', 100, 8500, 60, 'http://localhost:9000/lamp-images/8.jpg', false),
  ('Светильник встраиваемый светодиодный Forte 15W 4000K титан', 15, 1275, 120, 'http://localhost:9000/lamp-images/9.jpg', false),
  ('Подвесной светильник со стеклянными плафонами', 40, 3200, 120, 'http://localhost:9000/lamp-images/10.jpg', false),
  ('Светильник потолочный светодиодный Tend 9W 4000K черный', 9, 765, 120, 'http://localhost:9000/lamp-images/11.jpg', false),
  ('Подвесной светодиодный светильник', 40, 3200, 120, 'http://localhost:9000/lamp-images/12.jpg', false);

INSERT INTO users (username, passwrd, is_mod) VALUES
  ('test_user', 'test_password', false);

INSERT INTO light_requests (user_id, moderator_id, max_total_power_w) VALUES
  (1, NULL, 200);

INSERT INTO light_request_to_lamps (request_id, lamp_id, area_m2, number) VALUES
  (1, 2, 40, 1),
  (1, 3, 25, 10);
