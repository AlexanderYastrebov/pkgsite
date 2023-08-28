var s=class{constructor(){this.expand=document.querySelector(".js-versionsExpand");this.collapse=document.querySelector(".js-versionsCollapse");this.details=[...document.querySelectorAll(".js-versionDetails")];if(this.expand?.parentElement){this.details.some(e=>e.tagName==="DETAILS")&&(this.expand.parentElement.style.display="block");for(let e of this.details)e.addEventListener("click",()=>{this.updateButtons()});this.expand?.addEventListener("click",()=>{this.details.map(e=>e.open=!0),this.updateButtons()}),this.collapse?.addEventListener("click",()=>{this.details.map(e=>e.open=!1),this.updateButtons()}),this.updateButtons(),this.setCurrent()}}setCurrent(){let e=document.querySelector(".js-canonicalURLPath")?.dataset?.canonicalUrlPath,t=document.querySelector(`.js-versionLink[href="${e}"]`);t&&(t.style.fontWeight="bold")}updateButtons(){setTimeout(()=>{if(!this.expand||!this.collapse)return;let e,t;for(let n of this.details)e=e||n.open,t=t||!n.open;this.expand.style.display=t?"inline-block":"none",this.collapse.style.display=t?"none":"inline-block"})}};new s;export{s as VersionsController};
/*!
 * @license
 * Copyright 2021 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
//# sourceMappingURL=versions.js.map
