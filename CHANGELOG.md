### v1.1.0
 - Updated Carbone default API verison to 4
 - Increase default timeout to 60 seconds
 - Added method `SetAPIHeaders` to add custom headers, headers will be added automatically to HTTP requests, usage:
  ```go
  csdk.SetAPIHeaders(map[string]string{
    "carbone-template-delete-after": "86400",
    "carbone-webhook-url": "https://...",
  })
  ```

### v1.0.4
 - Fix template automatic reupload

### v1.0.3
 - Update Carbone default API Version to 3

### v1.0.2
 - Add license and update repository URL

### v1.0.1
 - Fix the `Render` method, add test and comments

### v1.0.0
  - Release July 3rd, 2020
  - It is possible to interact with the Carbone Render API with the following methods:
    - AddTemplate: upload a template and return a templateID
    - GetTemplate: return an uploaded template from a templateID
    - DeleteTemplate: delete a template from a templateID
    - Render: render a report from a templateID
    - GenerateTemplateID: Pre compute the templateID
