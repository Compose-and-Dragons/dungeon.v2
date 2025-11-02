---
marp: true
html: true
theme: default
paginate: true
---
<style>
.dodgerblue {
  color: dodgerblue;
}
</style>
#### 🔍🔄 Tool Calls Detection Loop: [/compose-dragons/agents/npcagents.go#L493](/compose-dragons/agents/npcagents.go#L493)
```bash
	┌──────────────────────────────┐
	│DetectAndExecuteToolCalls     │
	│WithConfirmation              │
	└──────────────┬───────────────┘
	               │
	               │ Passes executeToolWithConfirmation
	               │ as executor
	               │
	               ▼
	┌──────────────────────────────┐
	│detectAndExecuteToolCallsLoop │
	└──────────────┬───────────────┘
	               │
	               │ For each tool request,
	               │ calls executor
	               │
	               ▼
	┌──────────────────────────────┐
	│executeToolWithConfirmation   │
	└──────────────┬───────────────┘
	               │
	               │ Ask user: y/n/q?
	               │
	    ┌──────────┼──────────┐
	    │          │          │
	    ▼          ▼          ▼
	 ┌───┐      ┌───┐      ┌───┐
	 │ y │      │ n │      │ q │
	 └─┬─┘      └─┬─┘      └─┬─┘
	   │          │          │
	   │          │          └──► Set stopped=true
	   │          │
	   │          └──► Append "cancelled" to history
	   │
	   └──► Call executeTool() (Append result to history)
```