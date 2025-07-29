# KScript

- run a script by defining map data.
- run a script file by input arguments.

## DEV

```text
   ScriptRunner
        |
        | scan and load script config files. eg: yaml
        |
   ScriptEntry map
        |
   -----|---------
   |		     |
   |		     |
Script-Task    Script-File
   |               |
Commands        Script file info
   |               |
   |               |
    \           /
   Run by input and context
```

```mermaid
graph TD
    subgraph 初始化流程
        A[ScriptRunner] --> B[扫描并加载脚本配置文件]
        B --> C{解析配置文件}
        C -->|YAML/JSON等格式| D[生成 ScriptEntry 映射表]
    end

    subgraph 映射表处理
        D --类型:Define--> E[Define-Map 处理]
        D --类型:File--> F[Script-File 处理]

        E --> E1[注册命令定义]
        E1 --> E2[绑定命令参数]

        F --> F1[解析脚本文件元数据]
        F1 --> F2[验证文件完整性]
    end

    subgraph 执行阶段
        E2 --> G[执行引擎]
        F2 --> G
        G --> H[根据输入参数和运行时上下文执行]
        H --> I[返回执行结果]
    end

    classDef default fill:#f8f9fa,stroke:#333
    classDef process fill:#d1e7dd,stroke:#155724
    classDef entity fill:#fff3cd,stroke:#856404

    class A,B,C,D process
    class E,F,E1,E2,F1,F2 entity
    class G,H,I default
```
