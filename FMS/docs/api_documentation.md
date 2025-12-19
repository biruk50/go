# FMS API Documentation

This document summarizes core endpoints and RBAC.

## Auth

- POST /register -> register user
- POST /login -> obtain JWT

## Users (protected)

- GET /users -> list users (Finance/Admin)
- GET /users/me -> current user
- PUT /users/:id/role -> update role (Finance only)

## Budgets

- POST /budgets -> create budget (General Staff)
- GET /budgets -> list budgets (protected)
- GET /budgets/:id -> detail
- PATCH/PUT /budgets/:id -> update before approval (owner)
- POST /budgets/:id/approve -> approve (Finance only)
- POST /budgets/:id/reject -> reject (Finance only)
- GET /budgets/:id/summary -> summary of usage

## Cash Requests

- POST /cash-requests -> submit request
- GET /cash-requests -> list
- GET /cash-requests/:id -> detail
- POST /cash-requests/:id/approve -> approve (Finance only)
- POST /cash-requests/:id/reject -> reject (Finance only)
- POST /cash-requests/:id/disburse -> disburse funds (Finance only)

## Expenses

- POST /expenses -> record expense
- GET /expenses -> list
- GET /expenses/:id -> detail
- POST /expenses/:id/receipts -> attach receipt (uploads accepted as URL)
- PUT /expenses/:id/verify -> mark verified (Finance only)

## Reports

- GET /reports/overview -> high level overview (Finance)
- GET /reports/budgets -> budgets report
- GET /reports/cash-requests -> cash requests report
- GET /reports/expenses -> expense report

## RBAC Notes

- Finance role: approve/reject/disburse/verify and full report access.
- General Staff: submit budgets, cash requests, expenses and view own data.
