# Reflection: Reflect+Archive (повторная)

**Date:** 2025-03-01

## Контекст

Повторная процедура reflect+archive. tasks.md пуст.

## Что сделано после предыдущей архивации

- Transaction Update: использование old.AccountId и old.Tokens при изменении баланса; автоустановка date_approved/date_cancelled при смене статуса.
- TransactionLogic.Update: убран лишний GetByID.
- CreateWithReserveUpdate: блокировка строки аккаунта (SELECT FOR UPDATE) перед SUM.
- Константы: удалён TxStatusRefunded; добавлен PaymentTypeBonus (3).
- SQL: status check (0–3), payment_type (0–3), комментарии.
- Proto: обновлены комментарии статусов и PaymentType (BONUS), убран REFUNDED.

## Результат

tasks.md пуст; созданы reflection-20250301-s2.md и archive-20250301-s2.md.
