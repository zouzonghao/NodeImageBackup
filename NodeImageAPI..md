# NodeImage API 文档示例

本文件为 NodeImage API 的使用示例，敏感信息（如 API Token）请使用您自己的密钥替换。

---

## 1. 上传图片

- **接口地址**：`POST https://api.nodeimage.com/api/upload`
- **请求头**：
  - `X-API-Key: YOUR_API_TOKEN_HERE`
- **请求体**：
  - `image`：待上传图片文件（multipart/form-data）

**请求示例：**
```bash
curl -X POST "https://api.nodeimage.com/api/upload" \
  -H "X-API-Key: YOUR_API_TOKEN_HERE" \
  -F "image=@/path/to/your/image.jpg"
```

**返回示例：**
```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "image_id": "Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J",
  "filename": "Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J.avif",
  "size": 25376,
  "links": {
    "direct": "https://cdn.nodeimage.com/i/Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J.avif",
    "html": "<img src=\"https://cdn.nodeimage.com/i/Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J.avif\" alt=\"image\">",
    "markdown": "![image](https://cdn.nodeimage.com/i/Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J.avif)",
    "bbcode": "[img]https://cdn.nodeimage.com/i/Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J.avif[/img]"
  }
}
```

---

## 2. 删除图片

- **接口地址**：`DELETE https://api.nodeimage.com/api/v1/delete/{image_id}`
- **请求头**：
  - `X-API-Key: YOUR_API_TOKEN_HERE`
- **路径参数**：
  - `image_id`：要删除的图片 ID

**请求示例：**
```bash
curl -X DELETE "https://api.nodeimage.com/api/v1/delete/{image_id}" \
  -H "X-API-Key: YOUR_API_TOKEN_HERE"
```

**返回示例：**
```json
{
  "success": true,
  "message": "删除成功"
}
```

---

## 3. 获取图片列表

- **接口地址**：`GET https://api.nodeimage.com/api/v1/list`
- **请求头**：
  - `X-API-Key: YOUR_API_TOKEN_HERE`

**请求示例：**
```bash
curl -X GET "https://api.nodeimage.com/api/v1/list" \
  -H "X-API-Key: YOUR_API_TOKEN_HERE"
```

**返回示例（部分字段）：**
```json
{
  "success": true,
  "images": [
    {
      "image_id": "Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J",
      "filename": "Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J.avif",
      "size": 25376,
      "uploaded_at": "2024-06-01T12:34:56Z",
      "links": {
        "direct": "https://cdn.nodeimage.com/i/Nx1mskpFq8BTSQVBGEnrDHSxnw95SH3J.avif"
      }
    }
    // ... 更多图片对象
  ]
}
```

---

如需更多接口说明，请参考官方文档或联系 API 提供方。 