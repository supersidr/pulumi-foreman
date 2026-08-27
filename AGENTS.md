# pulumi-foreman

Локальный Pulumi-провайдер поверх Terraform (`upstream/` + патчи). Правки upstream в GitHub не пушить.

## Патчи

- Изменения terraform-провайдера — только файлы в `patches/` (`0001-…`, `0010-…`).
- `upstream/` — submodule. Некоммиченный WIP host/subnet не смешивать с новым lookup-патчем.
- Новый data source / resource: API-клиент + resource + data source + тесты по образцу `environment`, регистрация в `upstream/foreman/provider.go`.
- В `provider/resources.go` для нового data source omit `__meta_` / `__meta__` (ломает codegen).

## Сборка (локально)

Версию передавать явно, не полагаться на `PROVIDER_VERSION ?= 1.0.1` в Makefile:

```bash
make schema generate_python provider PROVIDER_VERSION=x.y.z
```

- Не гонять `make generate` (dotnet/go/nodejs) без явного запроса.
- После tfgen не откатывать чужие поля host/subnet в Python SDK, если они уже в HEAD; новые файлы lookup оставить.
- Версия должна совпасть в `sdk/python/pyproject.toml`, `sdk/python/pulumi_foreman/pulumi-plugin.json` и в `FOREMAN_PROVIDER_VERSION` у потребителей (pulumi-modules).

Плагин:

```bash
pulumi plugin install resource foreman x.y.z --file bin/pulumi-resource-foreman
```

Python SDK ставить в то окружение, из которого запускается Pulumi (`pip install --force-reinstall --no-deps sdk/python`).

## Lookup по имени

Уже есть: `get_domain`, `get_smartproxy`, `get_hostgroup`, `get_environment` (puppet environment).

Location / Organization — `get_location` / `get_organization` (патч `0010-add-location-organization-lookup.patch`). Поиск без taxonomy (как hostgroup): провайдер без `location_id` / `organization_id`.

## Чего не делать

- Не коммитить и не пушить, пока не попросили.
- Не переинициализировать `upstream` (`scripts/upstream.sh init -f`), если в working tree есть нужный WIP.
- Не читать `sdk/dotnet`, `sdk/nodejs`, `sdk/go`, если задача только Python / бинарь.
