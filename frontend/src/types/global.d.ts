/**
 * 全域型別宣告
 */

// 允許 import CSS 檔案
declare module "*.css" {
  const content: { [className: string]: string };
  export default content;
}

// 允許 import 圖片檔案
declare module "*.svg" {
  const content: string;
  export default content;
}

declare module "*.png" {
  const content: string;
  export default content;
}

declare module "*.jpg" {
  const content: string;
  export default content;
}

declare module "*.jpeg" {
  const content: string;
  export default content;
}

declare module "*.gif" {
  const content: string;
  export default content;
}

declare module "*.webp" {
  const content: string;
  export default content;
}

