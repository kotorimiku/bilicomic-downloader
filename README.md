# bilicomic-downloader

一个基于 Wails + Vue3 + TypeScript 开发的哔哩哔哩漫画下载器，提供简洁易用的界面和丰富的导出选项。

## ✨ 主要特性

- 🚀 **批量下载**：支持哔哩哔哩漫画章节批量下载
- 📚 **多格式导出**：支持 EPUB、ZIP、图片等多种导出格式
- 🎨 **界面友好**：基于 Vue3 构建的现代化用户界面
- ⚙️ **自定义配置**：支持自定义命名规则、图片格式等设置
- 📄 **元数据支持**：自动生成 ComicInfo.xml 元数据文件

## 🛠️ 前置条件

- Go 1.18 或更高版本
- Node.js 16 或更高版本
- pnpm（推荐）或 npm/yarn
- Wails v2

## 🚀 快速开始

### 克隆项目

```bash
git clone https://github.com/your-username/bilicomic-downloader.git
cd bilicomic-downloader
```

### 安装依赖

```bash
# 安装前端依赖
cd frontend
pnpm install
cd ..

# 安装 Wails（如果未安装）
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 开发模式

```bash
wails dev
```

应用将在开发模式下启动，支持热重载。

### 构建发布版

```bash
wails build
```

构建完成后，可执行文件将生成在 `build/bin/` 目录中。

## ⚙️ 配置

- **项目配置**：编辑 `wails.json` 文件
- **用户设置**：在应用界面的"设置"选项和`bcconfig.json` 文件中进行个性化配置

## 📖 使用说明

1. 启动应用
2. 在界面中输入哔哩哔哩漫画链接
3. 选择需要下载的章节
4. 配置导出格式和相关选项
5. 点击下载按钮开始下载

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进项目。

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详细信息。

## ⚠️ 免责声明

本工具仅供学习和个人使用，请遵守相关法律法规和平台服务条款。下载的内容版权归原作者所有，请勿用于商业用途。
