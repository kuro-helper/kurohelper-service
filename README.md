# kurohelper-service

KuroHelper 的共用 Go service module。它不持有 Discord Bot Token，也不直接註冊 Discord event handler。

## `kuro` package

- `client.go`：連接 `kurohelper-ai-runtime` 的 authenticated WebSocket client，包含斷線重連、request correlation 與 timeout。
- `types.go`：protocol v1 的生成、健康狀態與記憶操作型別。
- `access.go`：ID allowlist 與 Kuro 前綴/mention 觸發的純函式。
- `context.go`：把 Discord 最近訊息限制在指定數量/字元內，產生 prompt context、記憶檢索文字與去重後的事件參與者。

`typing` 與 `images` 是中間事件，不會提早結束等待中的文字生成；只有正式 response、error、逾時或斷線會完成 request。

## 資料層

`db/kuro_channel_state.go` 保存每個 Discord 頻道的最近 `小黑 /newchat` 邊界。這使 Bot 重啟後仍能維持頻道短期上下文的切割點。

## 使用原則

Discord API 與 UI handler 留在 `kurohelper`；可測試、可重用的協定、上下文與資料操作放在這個 module；SillyTavern/模型/記憶實作留在 `kurohelper-ai-runtime`。

```powershell
go test ./...
```
