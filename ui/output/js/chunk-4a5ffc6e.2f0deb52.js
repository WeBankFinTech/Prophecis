(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-4a5ffc6e"],{"1fe5":function(t,e,a){"use strict";a.r(e);var i=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("breadcrumb-nav",{staticClass:"margin-b20"},[a("el-button",{attrs:{type:"primary",icon:"el-icon-plus",size:"small"},on:{click:function(e){t.dialogVisible=!0}}},[t._v(" "+t._s(t.$t("data.addDataStore"))+" ")])],1),a("el-table",{staticStyle:{width:"100%"},attrs:{data:t.dataList,border:"","empty-text":t.$t("common.noData")}},[a("el-table-column",{attrs:{prop:"path",label:t.$t("data.dataStorePath")}}),a("el-table-column",{attrs:{prop:"isActive",label:t.$t("user.isActive"),formatter:t.translationActive}}),a("el-table-column",{attrs:{prop:"remarks",label:t.$t("DI.description")}}),a("el-table-column",{attrs:{label:t.$t("common.operation")},scopedSlots:t._u([{key:"default",fn:function(e){return[a("el-button",{attrs:{type:"text",size:"small"},on:{click:function(a){return t.updateItem(e.row)}}},[t._v(t._s(t.$t("common.modify")))]),a("el-button",{attrs:{type:"text",size:"small"},on:{click:function(a){return t.deleteItem(e.row)}}},[t._v(t._s(t.$t("common.delete")))]),a("el-button",{attrs:{type:"text",size:"small"},on:{click:function(a){return t.settingGroup(e.row)}}},[t._v(t._s(t.$t("user.userGroupSettings")))])]}}])})],1),a("el-pagination",{attrs:{"current-page":t.pagination.pageNumber,"page-sizes":t.sizeList,"page-size":t.pagination.pageSize,layout:"total, sizes, prev, pager, next, jumper",total:t.pagination.totalPage,background:""},on:{"size-change":t.handleSizeChange,"current-change":t.handleCurrentChange}}),a("el-dialog",{attrs:{title:t.title,visible:t.dialogVisible,"custom-class":"dialog-style",width:"800px"},on:{"update:visible":function(e){t.dialogVisible=e},close:t.clearDialogForm}},[a("el-form",{ref:"formValidate",attrs:{model:t.form,rules:t.ruleValidate,"label-width":"175px"}},[a("el-form-item",{attrs:{label:t.$t("data.dataStorePath"),prop:"path"}},[a("el-input",{attrs:{placeholder:t.$t("data.dataStorePro"),disabled:t.modify},model:{value:t.form.path,callback:function(e){t.$set(t.form,"path",e)},expression:"form.path"}})],1),a("el-form-item",{attrs:{label:t.$t("user.isActive"),prop:"isActive"}},[a("el-select",{attrs:{placeholder:t.$t("user.activePro")},model:{value:t.form.isActive,callback:function(e){t.$set(t.form,"isActive",e)},expression:"form.isActive"}},[a("el-option",{attrs:{label:t.$t("common.yes"),value:1}})],1)],1),a("el-form-item",{attrs:{label:t.$t("data.describe"),prop:"remarks"}},[a("el-input",{attrs:{type:"textarea",placeholder:t.$t("data.describePro")},model:{value:t.form.remarks,callback:function(e){t.$set(t.form,"remarks",e)},expression:"form.remarks"}})],1)],1),a("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[a("el-button",{attrs:{type:"primary",disabled:t.btnDisabled},on:{click:t.subInfo}},[t._v(t._s(t.$t("common.save")))]),a("el-button",{on:{click:function(e){t.dialogVisible=!1}}},[t._v(t._s(t.$t("common.cancel")))])],1)],1)],1)},o=[],r=(a("4d63"),a("ac1f"),a("25f0"),{data:function(){return{dialogVisible:!1,form:{remarks:"",isActive:1,path:""},modify:!1,dataList:[]}},mounted:function(){this.getListData()},computed:{title:function(){return this.modify?this.$t("data.modifyDataStore"):this.$t("data.newDataStore")},getUrl:function(){return"/cc/".concat(this.FesEnv.ccApiVersion,"/storages")},ruleValidate:function(){return this.resetFormFields(),{path:[{required:!0,message:this.$t("data.pathReq")},{pattern:new RegExp(/^\/[0-9a-zA-Z-_/*]*$/),message:this.$t("data.pathStartReq")}],isActive:[{required:!0,message:this.$t("user.activeReq")}],remarks:[{pattern:new RegExp(/^.{0,200}$/),message:this.$t("data.remarksMaxLength")}]}}},methods:{translationActive:function(t,e){return this.$t("common.yes")},updateItem:function(t){var e=this;this.dialogVisible=!0,this.modify=!0,setTimeout((function(){Object.assign(e.form,t)}),10)},deleteItem:function(t){var e="/cc/".concat(this.FesEnv.ccApiVersion,"/storages/path");this.deleteListItem(e,{path:t.path})},settingGroup:function(t){this.$router.push({name:"settingDataGroup",params:{storageId:t.id+""},query:{path:t.path}})},subInfo:function(){var t=this;this.$refs.formValidate.validate((function(e){if(e){t.setBtnDisabeld();var a=t.modify?"put":"post";t.FesApi.fetch("/cc/".concat(t.FesEnv.ccApiVersion,"/storages"),t.form,a).then((function(){t.getListData(),t.toast(),t.dialogVisible=!1}))}}))}}}),s=r,n=a("2877"),l=Object(n["a"])(s,i,o,!1,null,null,null);e["default"]=l.exports},"4d63":function(t,e,a){var i=a("83ab"),o=a("da84"),r=a("94ca"),s=a("7156"),n=a("9bf2").f,l=a("241c").f,c=a("44e7"),u=a("ad6d"),p=a("9f7f"),d=a("6eeb"),m=a("d039"),f=a("69f3").set,h=a("2626"),b=a("b622"),g=b("match"),v=o.RegExp,$=v.prototype,y=/a/g,k=/a/g,x=new v(y)!==y,w=p.UNSUPPORTED_Y,_=i&&r("RegExp",!x||w||m((function(){return k[g]=!1,v(y)!=y||v(k)==k||"/a/i"!=v(y,"i")})));if(_){var A=function(t,e){var a,i=this instanceof A,o=c(t),r=void 0===e;if(!i&&o&&t.constructor===A&&r)return t;x?o&&!r&&(t=t.source):t instanceof A&&(r&&(e=u.call(t)),t=t.source),w&&(a=!!e&&e.indexOf("y")>-1,a&&(e=e.replace(/y/g,"")));var n=s(x?new v(t,e):v(t,e),i?this:$,A);return w&&a&&f(n,{sticky:a}),n},V=function(t){t in A||n(A,t,{configurable:!0,get:function(){return v[t]},set:function(e){v[t]=e}})},S=l(v),z=0;while(S.length>z)V(S[z++]);$.constructor=A,A.prototype=$,d(o,"RegExp",A)}h("RegExp")},7156:function(t,e,a){var i=a("861d"),o=a("d2bb");t.exports=function(t,e,a){var r,s;return o&&"function"==typeof(r=e.constructor)&&r!==a&&i(s=r.prototype)&&s!==a.prototype&&o(t,s),t}}}]);
//# sourceMappingURL=chunk-4a5ffc6e.2f0deb52.js.map