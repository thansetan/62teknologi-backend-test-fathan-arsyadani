-- ----------------------------
-- Records of transactions
-- ----------------------------
INSERT INTO "public"."transactions" VALUES (1, 'delivery');
INSERT INTO "public"."transactions" VALUES (2, 'pickup');
INSERT INTO "public"."transactions" VALUES (3, 'restaurant_reservation');
SELECT setval(pg_get_serial_sequence('transactions', 'id'), coalesce(max(id), 0) + 1, false) FROM "transactions"; -- set auto increment val to max(id)+1