(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-1ce1c128"],{"0353":function(e,t,a){"use strict";a("3b23")},"09ca":function(e,t,a){},"21f4":function(e,t,a){"use strict";a("a9e3"),a("4d63"),a("ac1f"),a("25f0"),a("5319");var r={transDate:function(e,t){if(!e)return"";t=t||"yyyy-MM-dd hh:mm:ss",e=/^[0-9]*$/.test(e)?Number(e):new Date(e).getTime();var a=new Date(e),r={"y+":a.getFullYear(),"M+":a.getMonth()+1,"d+":a.getDate(),"h+":a.getHours(),"m+":a.getMinutes(),"s+":a.getSeconds()};for(var o in new RegExp("(y+)").test(t)&&(t=t.replace(RegExp.$1,r["y+"])),r)new RegExp("("+o+")").test(t)&&(t=t.replace(RegExp.$1,1===RegExp.$1.length?r[o]:("00"+r[o]).substr((""+r[o]).length)));return t},getTime:function(){var e=new Date,t=e.getFullYear()+"",a=e.getMonth()+1;a=a>9?a:"0"+a;var r=e.getDate();return r=r>9?r:"0"+r,t+a+r}};t["a"]=r},"232d":function(e,t,a){"use strict";a("09ca")},2626:function(e,t,a){"use strict";var r=a("d066"),o=a("9bf2"),s=a("b622"),n=a("83ab"),i=s("species");e.exports=function(e){var t=r(e),a=o.f;n&&t&&!t[i]&&a(t,i,{configurable:!0,get:function(){return this}})}},3142:function(e,t,a){},"3b23":function(e,t,a){},"4d63":function(e,t,a){var r=a("83ab"),o=a("da84"),s=a("94ca"),n=a("7156"),i=a("9bf2").f,l=a("241c").f,c=a("44e7"),u=a("ad6d"),m=a("9f7f"),p=a("6eeb"),d=a("d039"),f=a("69f3").set,g=a("2626"),h=a("b622"),b=h("match"),y=o.RegExp,v=y.prototype,I=/a/g,$=/a/g,x=new y(I)!==I,D=m.UNSUPPORTED_Y,E=r&&s("RegExp",!x||D||d((function(){return $[b]=!1,y(I)!=I||y($)==$||"/a/i"!=y(I,"i")})));if(E){var w=function(e,t){var a,r=this instanceof w,o=c(e),s=void 0===t;if(!r&&o&&e.constructor===w&&s)return e;x?o&&!s&&(e=e.source):e instanceof w&&(s&&(t=u.call(e)),e=e.source),D&&(a=!!t&&t.indexOf("y")>-1,a&&(t=t.replace(/y/g,"")));var i=n(x?new y(e,t):y(e,t),r?this:v,w);return D&&a&&f(i,{sticky:a}),i},P=function(e){e in w||i(w,e,{configurable:!0,get:function(){return y[e]},set:function(t){y[e]=t}})},k=l(y),_=0;while(k.length>_)P(k[_++]);v.constructor=w,w.prototype=v,p(o,"RegExp",w)}g("RegExp")},"4d7f":function(e,t,a){"use strict";a("3142")},"4de4":function(e,t,a){"use strict";var r=a("23e7"),o=a("b727").filter,s=a("1dde"),n=a("ae40"),i=s("filter"),l=n("filter");r({target:"Array",proto:!0,forced:!i||!l},{filter:function(e){return o(this,e,arguments.length>1?arguments[1]:void 0)}})},5319:function(e,t,a){"use strict";var r=a("d784"),o=a("825a"),s=a("7b0b"),n=a("50c4"),i=a("a691"),l=a("1d80"),c=a("8aa5"),u=a("14c3"),m=Math.max,p=Math.min,d=Math.floor,f=/\$([$&'`]|\d\d?|<[^>]*>)/g,g=/\$([$&'`]|\d\d?)/g,h=function(e){return void 0===e?e:String(e)};r("replace",2,(function(e,t,a,r){var b=r.REGEXP_REPLACE_SUBSTITUTES_UNDEFINED_CAPTURE,y=r.REPLACE_KEEPS_$0,v=b?"$":"$0";return[function(a,r){var o=l(this),s=void 0==a?void 0:a[e];return void 0!==s?s.call(a,o,r):t.call(String(o),a,r)},function(e,r){if(!b&&y||"string"===typeof r&&-1===r.indexOf(v)){var s=a(t,e,this,r);if(s.done)return s.value}var l=o(e),d=String(this),f="function"===typeof r;f||(r=String(r));var g=l.global;if(g){var $=l.unicode;l.lastIndex=0}var x=[];while(1){var D=u(l,d);if(null===D)break;if(x.push(D),!g)break;var E=String(D[0]);""===E&&(l.lastIndex=c(d,n(l.lastIndex),$))}for(var w="",P=0,k=0;k<x.length;k++){D=x[k];for(var _=String(D[0]),A=m(p(i(D.index),d.length),0),M=[],V=1;V<D.length;V++)M.push(h(D[V]));var O=D.groups;if(f){var S=[_].concat(M,A,d);void 0!==O&&S.push(O);var R=String(r.apply(void 0,S))}else R=I(_,d,A,M,O,r);A>=P&&(w+=d.slice(P,A)+R,P=A+_.length)}return w+d.slice(P)}];function I(e,a,r,o,n,i){var l=r+e.length,c=o.length,u=g;return void 0!==n&&(n=s(n),u=f),t.call(i,u,(function(t,s){var i;switch(s.charAt(0)){case"$":return"$";case"&":return e;case"`":return a.slice(0,r);case"'":return a.slice(l);case"<":i=n[s.slice(1,-1)];break;default:var u=+s;if(0===u)return t;if(u>c){var m=d(u/10);return 0===m?t:m<=c?void 0===o[m-1]?s.charAt(1):o[m-1]+s.charAt(1):t}i=o[u-1]}return void 0===i?"":i}))}}))},5530:function(e,t,a){"use strict";a.d(t,"a",(function(){return s}));a("a4d3"),a("4de4"),a("4160"),a("e439"),a("dbb4"),a("b64b"),a("159b");function r(e,t,a){return t in e?Object.defineProperty(e,t,{value:a,enumerable:!0,configurable:!0,writable:!0}):e[t]=a,e}function o(e,t){var a=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),a.push.apply(a,r)}return a}function s(e){for(var t=1;t<arguments.length;t++){var a=null!=arguments[t]?arguments[t]:{};t%2?o(Object(a),!0).forEach((function(t){r(e,t,a[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(a)):o(Object(a)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(a,t))}))}return e}},7156:function(e,t,a){var r=a("861d"),o=a("d2bb");e.exports=function(e,t,a){var s,n;return o&&"function"==typeof(s=t.constructor)&&s!==a&&r(n=s.prototype)&&n!==a.prototype&&o(e,n),e}},"99af":function(e,t,a){"use strict";var r=a("23e7"),o=a("d039"),s=a("e8b5"),n=a("861d"),i=a("7b0b"),l=a("50c4"),c=a("8418"),u=a("65f0"),m=a("1dde"),p=a("b622"),d=a("2d00"),f=p("isConcatSpreadable"),g=9007199254740991,h="Maximum allowed index exceeded",b=d>=51||!o((function(){var e=[];return e[f]=!1,e.concat()[0]!==e})),y=m("concat"),v=function(e){if(!n(e))return!1;var t=e[f];return void 0!==t?!!t:s(e)},I=!b||!y;r({target:"Array",proto:!0,forced:I},{concat:function(e){var t,a,r,o,s,n=i(this),m=u(n,0),p=0;for(t=-1,r=arguments.length;t<r;t++)if(s=-1===t?n:arguments[t],v(s)){if(o=l(s.length),p+o>g)throw TypeError(h);for(a=0;a<o;a++,p++)a in s&&c(m,p,s[a])}else{if(p>=g)throw TypeError(h);c(m,p++,s)}return m.length=p,m}})},a9e3:function(e,t,a){"use strict";var r=a("83ab"),o=a("da84"),s=a("94ca"),n=a("6eeb"),i=a("5135"),l=a("c6b6"),c=a("7156"),u=a("c04e"),m=a("d039"),p=a("7c73"),d=a("241c").f,f=a("06cf").f,g=a("9bf2").f,h=a("58a8").trim,b="Number",y=o[b],v=y.prototype,I=l(p(v))==b,$=function(e){var t,a,r,o,s,n,i,l,c=u(e,!1);if("string"==typeof c&&c.length>2)if(c=h(c),t=c.charCodeAt(0),43===t||45===t){if(a=c.charCodeAt(2),88===a||120===a)return NaN}else if(48===t){switch(c.charCodeAt(1)){case 66:case 98:r=2,o=49;break;case 79:case 111:r=8,o=55;break;default:return+c}for(s=c.slice(2),n=s.length,i=0;i<n;i++)if(l=s.charCodeAt(i),l<48||l>o)return NaN;return parseInt(s,r)}return+c};if(s(b,!y(" 0o1")||!y("0b1")||y("+0x1"))){for(var x,D=function(e){var t=arguments.length<1?0:e,a=this;return a instanceof D&&(I?m((function(){v.valueOf.call(a)})):l(a)!=b)?c(new y($(t)),a,D):$(t)},E=r?d(y):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,isFinite,isInteger,isNaN,isSafeInteger,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,parseFloat,parseInt,isInteger".split(","),w=0;E.length>w;w++)i(y,x=E[w])&&!i(D,x)&&g(D,x,f(y,x));D.prototype=v,v.constructor=D,n(o,b,D)}},b64b:function(e,t,a){var r=a("23e7"),o=a("7b0b"),s=a("df75"),n=a("d039"),i=n((function(){s(1)}));r({target:"Object",stat:!0,forced:i},{keys:function(e){return s(o(e))}})},cbb5:function(e,t,a){},dbb4:function(e,t,a){var r=a("23e7"),o=a("83ab"),s=a("56ef"),n=a("fc6a"),i=a("06cf"),l=a("8418");r({target:"Object",stat:!0,sham:!o},{getOwnPropertyDescriptors:function(e){var t,a,r=n(e),o=i.f,c=s(r),u={},m=0;while(c.length>m)a=o(r,t=c[m++]),void 0!==a&&l(u,t,a);return u}})},e439:function(e,t,a){var r=a("23e7"),o=a("d039"),s=a("fc6a"),n=a("06cf").f,i=a("83ab"),l=o((function(){n(1)})),c=!i||l;r({target:"Object",stat:!0,forced:c,sham:!i},{getOwnPropertyDescriptor:function(e,t){return n(s(e),t)}})},eefb:function(e,t,a){"use strict";a.r(t);var r=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",[a("breadcrumb-nav"),a("el-form",{ref:"queryValidate",staticClass:"query",attrs:{model:e.query,rules:e.ruleValidate,"label-width":"85px"}},[a("el-row",[a("el-col",{attrs:{span:6}},[a("el-form-item",{attrs:{label:e.$t("ns.nameSpace"),prop:"namespace"}},[a("el-select",{attrs:{filterable:"",placeholder:e.$t("ns.nameSpacePro")},model:{value:e.query.namespace,callback:function(t){e.$set(e.query,"namespace",t)},expression:"query.namespace"}},e._l(e.spaceOptionList,(function(e){return a("el-option",{key:e,attrs:{value:e,label:e}})})),1)],1)],1),a("el-col",{attrs:{span:6}},[a("el-form-item",{attrs:{label:e.$t("user.user"),prop:"userName"}},[a("el-select",{attrs:{filterable:"",clearable:"",placeholder:e.$t("user.userPro")},model:{value:e.query.userName,callback:function(t){e.$set(e.query,"userName",t)},expression:"query.userName"}},e._l(e.userOptionList,(function(e){return a("el-option",{key:e.id,attrs:{value:e.name,label:e.name}})})),1)],1)],1),a("el-col",{attrs:{span:12}},[a("el-button",{staticClass:"margin-l30",attrs:{type:"primary",icon:"el-icon-search"},on:{click:e.filterListData}},[e._v(" "+e._s(e.$t("common.filter"))+" ")]),a("el-button",{attrs:{type:"warning",icon:"el-icon-refresh-left"},on:{click:e.resetListData}},[e._v(" "+e._s(e.$t("common.reset"))+" ")]),a("el-button",{attrs:{type:"primary",icon:"plus-circle-o"},on:{click:e.showCreateDialog}},[e._v(" "+e._s(e.$t("AIDE.createNotebook"))+" ")])],1)],1)],1),a("el-table",{staticStyle:{width:"100%"},attrs:{data:e.dataList,border:"","empty-text":e.$t("common.noData")}},[a("el-table-column",{attrs:{prop:"name",label:e.$t("DI.name"),"min-width":"100"}}),a("el-table-column",{attrs:{prop:"user",label:e.$t("user.user")}}),a("el-table-column",{attrs:{prop:"namespace",label:e.$t("ns.nameSpace"),"min-width":"220"}}),a("el-table-column",{attrs:{prop:"cpu",label:e.$t("AIDE.cpu"),"min-width":"50"}}),a("el-table-column",{attrs:{prop:"gpu",label:e.$t("AIDE.gpu"),"min-width":"50"}}),a("el-table-column",{attrs:{prop:"memory",label:e.$t("AIDE.memory1"),"min-width":"70"}}),a("el-table-column",{attrs:{prop:"image",label:e.$t("DI.image"),"min-width":"300"}}),a("el-table-column",{attrs:{prop:"status",label:e.$t("DI.status")},scopedSlots:e._u([{key:"default",fn:function(t){return[a("status-tag",{attrs:{status:t.row.status,statusObj:e.statusObj}})]}}])}),a("el-table-column",{attrs:{prop:"uptime",label:e.$t("DI.createTime"),"min-width":"125"}}),a("el-table-column",{attrs:{label:e.$t("common.operation"),"min-width":"85"},scopedSlots:e._u([{key:"default",fn:function(t){return[a("el-dropdown",{attrs:{size:"medium"},on:{command:e.handleCommand}},[a("el-button",{staticClass:"multi-operate-btn",attrs:{type:"text",icon:"el-icon-more",round:""}}),a("el-dropdown-menu",{attrs:{slot:"dropdown"},slot:"dropdown"},[a("el-dropdown-item",{attrs:{command:e.beforeHandleCommand("visitNotebook",t.row),disabled:"Ready"!==t.row.status}},[e._v(e._s(e.$t("AIDE.access")))]),a("el-dropdown-item",{attrs:{command:e.beforeHandleCommand("deleteItem",t.row)}},[e._v(e._s(e.$t("common.delete")))]),a("el-dropdown-item",{attrs:{command:e.beforeHandleCommand("copyJob",t.row)}},[e._v(e._s(e.$t("common.copy")))]),a("el-dropdown-item",{attrs:{command:e.beforeHandleCommand("editYarn",t.row)}},[e._v(e._s(e.$t("AIDE.editYarn")))])],1)],1)]}}])})],1),a("el-pagination",{attrs:{"current-page":e.pagination.pageNumber,"page-sizes":e.sizeList,"page-size":e.pagination.pageSize,layout:"total, sizes, prev, pager, next, jumper",total:e.pagination.totalPage,background:""},on:{"size-change":e.handleSizeChange,"current-change":e.handleCurrentChange}}),a("el-dialog",{attrs:{title:e.$t("AIDE.createNotebook"),visible:e.dialogVisible,"custom-class":"dialog-style","close-on-click-modal":!1,width:"1150px"},on:{"update:visible":function(t){e.dialogVisible=t},close:e.clearDialogForm1}},[a("step",{attrs:{step:e.currentStep,"step-list":e.stepList}}),a("el-form",{ref:"formValidate",staticClass:"add-node-form",attrs:{model:e.form,rules:e.ruleValidate,"label-width":"155px"}},[0===e.currentStep?a("div",{key:"step0"},[a("div",{staticClass:"subtitle"},[e._v(" "+e._s(e.$t("DI.basicSettings"))+" ")]),a("el-row",[a("el-col",{attrs:{span:24}},[a("el-form-item",{attrs:{label:e.$t("AIDE.NotebookName"),prop:"name",maxlength:"225"}},[a("el-input",{attrs:{placeholder:e.$t("AIDE.NotebookPro")},model:{value:e.form.name,callback:function(t){e.$set(e.form,"name",t)},expression:"form.name"}})],1)],1)],1)],1):e._e(),1===e.currentStep?a("div",{key:"step1"},[a("div",{staticClass:"subtitle"},[e._v(" "+e._s(e.$t("DI.imageSettings"))+" ")]),a("el-row",[a("el-form-item",{attrs:{label:e.$t("DI.imageType"),prop:"image.imageType"}},[a("el-radio-group",{on:{change:e.changeImageType},model:{value:e.form.image.imageType,callback:function(t){e.$set(e.form.image,"imageType",t)},expression:"form.image.imageType"}},[a("el-radio",{attrs:{label:"Standard"}},[e._v(e._s(e.$t("DI.standard")))]),a("el-radio",{attrs:{label:"Custom"}},[e._v(e._s(e.$t("DI.custom")))])],1)],1)],1),a("el-row",[e.isCustom?e._e():a("el-form-item",{attrs:{label:e.$t("DI.imageSelection"),prop:"image.imageOption"}},[a("span",[e._v(" "+e._s(e.defineImage)+" ")]),a("el-select",{attrs:{placeholder:e.$t("DI.imagePro")},model:{value:e.form.image.imageOption,callback:function(t){e.$set(e.form.image,"imageOption",t)},expression:"form.image.imageOption"}},e._l(e.imageOptionList,(function(e,t){return a("el-option",{key:t,attrs:{label:e,value:e}})})),1)],1),e.isCustom?a("el-form-item",{key:"imageInput",attrs:{label:e.$t("DI.imageSelection"),prop:"image.imageInput"}},[a("span",[e._v(" "+e._s(e.defineImage)+" ")]),a("el-input",{attrs:{placeholder:e.$t("DI.imageInputPro")},model:{value:e.form.image.imageInput,callback:function(t){e.$set(e.form.image,"imageInput",t)},expression:"form.image.imageInput"}})],1):e._e()],1)],1):e._e(),2===e.currentStep?a("div",{key:"step2"},[a("div",{staticClass:"subtitle"},[e._v(" "+e._s(e.$t("DI.computingResource"))+" ")]),a("el-row",[a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("ns.nameSpace"),prop:"namespace"}},[a("el-select",{attrs:{filterable:"",placeholder:e.$t("ns.nameSpacePro")},model:{value:e.form.namespace,callback:function(t){e.$set(e.form,"namespace",t)},expression:"form.namespace"}},e._l(e.spaceOptionList,(function(e){return a("el-option",{key:e,attrs:{label:e,value:e}})})),1)],1)],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("AIDE.cpu"),prop:"cpu"}},[a("el-input",{attrs:{placeholder:e.$t("AIDE.CPUPro")},model:{value:e.form.cpu,callback:function(t){e.$set(e.form,"cpu",t)},expression:"form.cpu"}},[a("span",{attrs:{slot:"append"},slot:"append"},[e._v("Core")])])],1)],1)],1),a("el-row",[a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("AIDE.memory"),prop:"memory"}},[a("el-input",{attrs:{placeholder:e.$t("AIDE.memoryPro")},model:{value:e.form.memory,callback:function(t){e.$set(e.form,"memory",t)},expression:"form.memory"}},[a("span",{attrs:{slot:"append"},slot:"append"},[e._v("Gi")])])],1)],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("AIDE.gpu"),prop:"extraResources"}},[a("el-input",{attrs:{placeholder:e.$t("AIDE.gpuPro")},model:{value:e.form.extraResources,callback:function(t){e.$set(e.form,"extraResources",t)},expression:"form.extraResources"}},[a("span",{attrs:{slot:"append"},slot:"append"},[e._v(e._s(e.$t("AIDE.block")))])])],1)],1)],1),a("el-form-item",{attrs:{label:e.$t("AIDE.yarnClusterResources")}},[a("el-switch",{model:{value:e.form.cluster,callback:function(t){e.$set(e.form,"cluster",t)},expression:"form.cluster"}})],1),!0===e.form.cluster?a("el-row",[a("el-form-item",{attrs:{label:e.$t("DI.YARNQueue"),prop:"queue"}},[a("el-input",{attrs:{maxlength:"225",placeholder:e.$t("DI.queuePro")},model:{value:e.form.queue,callback:function(t){e.$set(e.form,"queue",t)},expression:"form.queue"}})],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("DI.driverMemory"),prop:"driverMemory"}},[a("el-input",{attrs:{maxlength:"10",placeholder:e.$t("DI.driverMemoryPro")},model:{value:e.form.driverMemory,callback:function(t){e.$set(e.form,"driverMemory",t)},expression:"form.driverMemory"}},[a("span",{attrs:{slot:"append"},slot:"append"},[e._v("Gi")])])],1)],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("DI.executorInstances"),prop:"executorCores"}},[a("el-input",{attrs:{maxlength:"10",placeholder:e.$t("DI.executorInstancesPro")},model:{value:e.form.executorCores,callback:function(t){e.$set(e.form,"executorCores",t)},expression:"form.executorCores"}},[a("span",{attrs:{slot:"append"},slot:"append"},[e._v("Core")])])],1)],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("DI.executorMemory"),prop:"executorMemory"}},[a("el-input",{attrs:{maxlength:"10",placeholder:e.$t("DI.executorMemoryPro")},model:{value:e.form.executorMemory,callback:function(t){e.$set(e.form,"executorMemory",t)},expression:"form.executorMemory"}},[a("span",{attrs:{slot:"append"},slot:"append"},[e._v("Gi")])])],1)],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("DI.linkisInstance"),prop:"executors"}},[a("el-input",{attrs:{maxlength:"10",placeholder:e.$t("DI.linkisInstancePro")},model:{value:e.form.executors,callback:function(t){e.$set(e.form,"executors",t)},expression:"form.executors"}})],1)],1)],1):e._e()],1):e._e(),3===e.currentStep?a("div",{key:"step3"},[a("el-row",[a("el-form-item",{attrs:{label:e.$t("AIDE.settingType")}},[a("el-radio-group",{model:{value:e.storageDefalut,callback:function(t){e.storageDefalut=t},expression:"storageDefalut"}},[a("el-radio",{attrs:{label:"true"}},[e._v(e._s(e.$t("AIDE.default")))]),a("el-radio",{attrs:{label:"false"}},[e._v(e._s(e.$t("DI.custom")))])],1)],1)],1),"true"==e.storageDefalut?a("div",[a("div",{staticClass:"subtitle"},[e._v(" "+e._s(e.$t("AIDE.workspace"))+" ")]),a("el-row",[a("el-form-item",{attrs:{label:e.$t("DI.localPath"),prop:"workspaceVolume.localPath"}},[a("el-select",{attrs:{filterable:"",placeholder:e.$t("AIDE.workLocalPathNodePro")},model:{value:e.form.workspaceVolume.localPath,callback:function(t){e.$set(e.form.workspaceVolume,"localPath",t)},expression:"form.workspaceVolume.localPath"}},e._l(e.localPathList,(function(e,t){return a("el-option",{key:t,attrs:{label:e,value:e}})})),1)],1),a("el-form-item",{attrs:{label:e.$t("AIDE.mountPath"),prop:"workspaceVolume.mountPath"}},[a("el-input",{attrs:{placeholder:e.$t("AIDE.mountPathPro"),disabled:""},model:{value:e.form.workspaceVolume.mountPath,callback:function(t){e.$set(e.form.workspaceVolume,"mountPath",t)},expression:"form.workspaceVolume.mountPath"}})],1)],1)],1):e._e(),"false"==e.storageDefalut?a("div",{staticClass:"subtitle data-space"},[e._v(" "+e._s(e.$t("AIDE.dataStorageSettings"))+" ")]):e._e(),"false"==e.storageDefalut?a("div",e._l(e.form.dataVolume,(function(t,r){return a("div",{key:"dataVolume"+r},[a("el-row",[a("el-form-item",{attrs:{label:e.$t("DI.localPath"),prop:"dataVolume["+r+"].localPath"}},[a("el-select",{attrs:{filterable:"",placeholder:e.$t("AIDE.dataLocalPathNodePro")},model:{value:t.localPath,callback:function(a){e.$set(t,"localPath",a)},expression:"item.localPath"}},e._l(e.localPathList,(function(e,t){return a("el-option",{key:t,attrs:{label:e,value:e}})})),1)],1),a("el-form-item",{attrs:{label:e.$t("AIDE.subpath"),prop:"dataVolume["+r+"].subPath"}},[a("el-input",{attrs:{placeholder:e.$t("AIDE.subpathPro")},model:{value:t.subPath,callback:function(a){e.$set(t,"subPath",a)},expression:"item.subPath"}})],1),a("el-form-item",{attrs:{label:e.$t("AIDE.mountPath"),prop:"dataVolume["+r+"].mountPath"}},[a("el-input",{attrs:{placeholder:e.$t("AIDE.dataMountPathPro")},model:{value:t.mountPath,callback:function(a){e.$set(t,"mountPath",a)},expression:"item.mountPath"}})],1)],1),a("el-row",[a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("AIDE.mountType"),prop:"dataVolume["+r+"].mountType"}},[a("el-select",{attrs:{placeholder:e.$t("AIDE.mountTypePro")},model:{value:t.mountType,callback:function(a){e.$set(t,"mountType",a)},expression:"item.mountType"}},[a("el-option",{attrs:{value:"New",label:e.$t("AIDE.nodeDirectory")}})],1)],1)],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("AIDE.accessMode"),prop:"dataVolume["+r+"].accessMode"}},[a("el-select",{attrs:{placeholder:e.$t("AIDE.accessModePro")},model:{value:t.accessMode,callback:function(a){e.$set(t,"accessMode",a)},expression:"item.accessMode"}},[a("el-option",{attrs:{value:"ReadOnlyMany",label:e.$t("AIDE.readOnly")}}),a("el-option",{attrs:{value:"ReadWriteMany",label:e.$t("AIDE.readWrite")}})],1)],1)],1)],1)],1)})),0):e._e()],1):e._e()]),a("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[e.currentStep>0?a("el-button",{attrs:{type:"primary"},on:{click:e.lastStep}},[e._v(" "+e._s(e.$t("common.lastStep"))+" ")]):e._e(),e.currentStep<3?a("el-button",{attrs:{type:"primary"},on:{click:e.nextStep}},[e._v(" "+e._s(e.$t("common.nextStep"))+" ")]):e._e(),3===e.currentStep?a("el-button",{attrs:{disabled:e.btnDisabled,type:"primary"},on:{click:e.subInfo}},[e._v(" "+e._s(e.$t("common.save"))+" ")]):e._e(),a("el-button",{on:{click:function(t){e.dialogVisible=!1}}},[e._v(" "+e._s(e.$t("common.cancel"))+" ")])],1)],1),a("el-dialog",{attrs:{title:e.$t("AIDE.yarnClusterResources"),visible:e.yarnDialog,"custom-class":"dialog-style","close-on-click-modal":!1,width:"1100px"},on:{"update:visible":function(t){e.yarnDialog=t}}},[a("el-form",{ref:"yarnForm",attrs:{rules:e.ruleValidate,model:e.yarnForm,"label-width":"190px"}},[a("el-row",[a("el-form-item",{attrs:{label:e.$t("DI.YARNQueue"),prop:"queue"}},[a("el-input",{attrs:{maxlength:"225",placeholder:e.$t("DI.queuePro")},model:{value:e.yarnForm.queue,callback:function(t){e.$set(e.yarnForm,"queue",t)},expression:"yarnForm.queue"}})],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("DI.driverMemory"),prop:"driverMemory"}},[a("el-input",{attrs:{maxlength:"10",placeholder:e.$t("DI.driverMemoryPro")},model:{value:e.yarnForm.driverMemory,callback:function(t){e.$set(e.yarnForm,"driverMemory",t)},expression:"yarnForm.driverMemory"}},[a("span",{attrs:{slot:"append"},slot:"append"},[e._v("Gi")])])],1)],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("DI.executorInstances"),prop:"executorCores"}},[a("el-input",{attrs:{maxlength:"10",placeholder:e.$t("DI.executorInstancesPro")},model:{value:e.yarnForm.executorCores,callback:function(t){e.$set(e.yarnForm,"executorCores",t)},expression:"yarnForm.executorCores"}},[a("span",{attrs:{slot:"append"},slot:"append"},[e._v("Core")])])],1)],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("DI.executorMemory"),prop:"executorMemory"}},[a("el-input",{attrs:{maxlength:"10",placeholder:e.$t("DI.executorMemoryPro")},model:{value:e.yarnForm.executorMemory,callback:function(t){e.$set(e.yarnForm,"executorMemory",t)},expression:"yarnForm.executorMemory"}},[a("span",{attrs:{slot:"append"},slot:"append"},[e._v("Gi")])])],1)],1),a("el-col",{attrs:{span:12}},[a("el-form-item",{attrs:{label:e.$t("DI.linkisInstance"),prop:"executors"}},[a("el-input",{attrs:{maxlength:"10",placeholder:e.$t("DI.linkisInstancePro")},model:{value:e.yarnForm.executors,callback:function(t){e.$set(e.yarnForm,"executors",t)},expression:"yarnForm.executors"}})],1)],1)],1)],1),a("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[a("el-button",{attrs:{type:"primary"},on:{click:e.saveYarnForm}},[e._v(" "+e._s(e.$t("common.save"))+" ")]),a("el-button",{on:{click:e.cannelSaveYarnForm}},[e._v(" "+e._s(e.$t("common.cancel"))+" ")])],1)],1)],1)},o=[],s=(a("99af"),a("4160"),a("c975"),a("b0c0"),a("a9e3"),a("dca8"),a("4d63"),a("ac1f"),a("25f0"),a("159b"),a("5530")),n=a("b85c"),i=a("21f4"),l=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"step-steps"},[a("ul",e._l(e.stepList,(function(t,r){return a("li",{key:r},[a("span",{class:{"steps-dot":"steps-dot","steps-current":r===e.step,"steps-finish":e.step>r}}),e._v(" "+e._s(t)+" ")])})),0)])},c=[],u={props:{step:{type:Number,default:function(){return 0}},stepList:{type:Array,default:function(){return[]}}}},m=u,p=(a("232d"),a("2877")),d=Object(p["a"])(m,l,c,!1,null,"283c0962",null),f=d.exports,g={methods:{getDIUserOption:function(){var e=this;this.FesApi.fetch("/cc/".concat(this.FesEnv.ccApiVersion,"/users/myUsers"),"get").then((function(t){e.userOptionList=e.formatResList(t)}))},getDISpaceOption:function(){var e=this;this.FesApi.fetch("/cc/".concat(this.FesEnv.ccApiVersion,"/namespaces/myNamespace"),{},"get").then((function(t){e.spaceOptionList=e.formatResList(t)}))},getDIStoragePath:function(){var e=this;this.FesApi.fetch("/cc/".concat(this.FesEnv.ccApiVersion,"/groups/group/storage"),{},"get").then((function(t){e.localPathList=e.formatResList(t)}))},resetQueryFields:function(){var e=this;setTimeout((function(){e.$refs.queryValidate&&e.$refs.queryValidate.resetFields()}),10)},filterListData:function(){var e=this;this.$refs.queryValidate.validate((function(t){t&&(e.pagination.pageNumber=1,e.getListData())}))},resetListData:function(){this.pagination.pageNumber=1,this.$refs.queryValidate.resetFields(),this.getListData()}}},h=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("el-tag",{attrs:{type:e.statusObj[e.status]}},[e._v(e._s(e.status))])},b=[],y=(a("cbb5"),a("450d"),a("8bbc")),v=a.n(y),I={components:{ElTag:v.a},props:{statusObj:{type:Object,default:function(){return{success:"success",info:"info",warning:"warning",danger:"danger"}}},status:{type:String,default:"success"}},data:function(){return{}}},$=I,x=(a("0353"),Object(p["a"])($,h,b,!1,null,"0d8a0704",null)),D=x.exports,E={mixins:[g],components:{Step:f,StatusTag:D},data:function(){return{dialogVisible:!1,query:{namespace:"",userName:""},currentStep:0,form:{name:"",cpu:"",extraResources:"",memory:"",queue:"",namespace:"",driverMemory:"",executorCores:"",executorMemory:"",executors:"",cluster:!1,image:{imageType:"Standard",imageOption:"",imageInput:""},workspaceVolume:{accessMode:"",localPath:"",subPath:"",mountPath:"",mountType:"",size:0},dataVolume:[{accessMode:"",localPath:"",subPath:"",mountPath:"",mountType:"",size:0}]},localPathList:[],loading:!1,dataList:[],userOptionList:[],spaceOptionList:[],yarnDialog:!1,yarnForm:{cluster:!0,queue:"",driverMemory:"",executorCores:"",executorMemory:"",executors:""},storageDefalut:"true",statusObj:{Ready:"success",Waiting:"",Terminate:"danger"},setIntervalList:null}},mounted:function(){this.getListData(),this.getDIUserOption(),this.getDISpaceOption(),this.getDIStoragePath(),this.setIntervalList=setInterval(this.getListData,6e4)},computed:{noDataText:function(){return this.loading?"":this.$t("common.noData")},stepList:function(){return[this.$t("DI.basicSettings"),this.$t("DI.imageSettings"),this.$t("DI.computingResource"),this.$t("AIDE.storageSettings")]},defineImage:function(){return this.FesEnv.AIDE.defineImage},imageOptionList:function(){return this.FesEnv.AIDE.imageOption},userId:function(){return localStorage.getItem("userId")},isCustom:function(){return"Custom"===this.form.image.imageType},ruleValidate:function(){return this.resetQueryFields(),this.resetFormFields(),{name:[{required:!0,message:this.$t("AIDE.NotebookNameReq")},{type:"string",pattern:new RegExp(/^[a-z][a-z0-9-]*$/),message:this.$t("AIDE.NotebookNameFormatReq")}],namespace:[{required:!0,message:this.$t("ns.nameSpaceReq")}],"image.imageType":[{required:!0,message:this.$t("DI.imageTypeReq")}],"image.imageOption":[{required:!0,message:this.$t("DI.imageColumnReq")}],"image.imageInput":[{required:!0,message:this.$t("DI.imageColumnReq")},{pattern:new RegExp(/^[a-zA-Z0-9][a-zA-Z0-9-._]*$/),message:this.$t("DI.imageInputFormat")}],cpu:[{required:!0,message:this.$t("DI.CPUReq")},{pattern:new RegExp(/^((0{1})([.]\d{1})|([1-9]\d*)([.]\d{1})?)$/),message:this.$t("AIDE.CPUNumberReq")}],memory:[{required:!0,message:this.$t("DI.memoryReq")},{pattern:new RegExp(/^[1-9]\d*([.]0)*$/),message:this.$t("AIDE.memoryNumberReq")}],queue:[{required:!0,message:this.$t("DI.queueReq")},{pattern:new RegExp(/^[a-zA-Z0-9][a-zA-Z0-9_.]*$/),message:this.$t("DI.queueFormat")}],driverMemory:[{pattern:new RegExp(/^[1-9]\d*([.]0)*$/),message:this.$t("DI.driverMemoryFormat")},{pattern:new RegExp(/^.{0,10}$/),message:this.$t("AIDE.driverMemoryMax")}],executorCores:[{pattern:new RegExp(/^[1-9]\d*([.]0)*$/),message:this.$t("DI.executorInstancesFormat")},{pattern:new RegExp(/^.{0,10}$/),message:this.$t("AIDE.executorInstancesMax")}],executors:[{pattern:new RegExp(/^[1-9]\d*([.]0)*$/),message:this.$t("AIDE.executorFormat")},{pattern:new RegExp(/^.{0,10}$/),message:this.$t("AIDE.executorsMax")}],executorMemory:[{pattern:new RegExp(/^[1-9]\d*([.]0)*$/),message:this.$t("DI.executorMemoryFormat")},{pattern:new RegExp(/^.{0,10}$/),message:this.$t("AIDE.executorMemoryMax")}],"workspaceVolume.localPath":[{required:!0,message:this.$t("DI.localPathReq")},{pattern:new RegExp(/^\/\w[0-9a-zA-Z-_/]*$/),message:this.$t("AIDE.localPathFormat")}],"workspaceVolume.mountPath":[{required:!0,message:this.$t("AIDE.mountPathReq")}],"dataVolume[0].localPath":[{required:!0,message:this.$t("DI.localPathReq")},{pattern:new RegExp(/^\/\w[0-9a-zA-Z-_/]*$/),message:this.$t("AIDE.localPathFormat")}],"dataVolume[0].subPath":[{pattern:new RegExp(/^[a-zA-Z][0-9a-zA-Z-_/]*$/),message:this.$t("AIDE.subpathReq")}],"dataVolume[0].mountPath":[{required:!0,message:this.$t("AIDE.mountPathReq")},{pattern:new RegExp(/^\/\w[0-9a-zA-Z-_/]*$/),message:this.$t("AIDE.dataMountPathFormat")}],"dataVolume[0].mountType":[{required:!0,message:this.$t("AIDE.mountTypeReq")}],"dataVolume[0].accessMode":[{required:!0,message:this.$t("AIDE.accessModeReq")}],extraResources:[{pattern:new RegExp(/^(([1-9]{1}\d*([.]0)*)|([0]{1}))$/),message:this.$t("AIDE.GPUFormat")}]}}},methods:{showCreateDialog:function(){this.form.workspaceVolume.mountPath="/home/".concat(this.userId,"/workspace"),this.dialogVisible=!0},visitNotebook:function(e){var t="/cc/".concat(this.FesEnv.ccApiVersion,"/auth/access/namespaces/").concat(e.namespace,"/notebooks/").concat(e.name),a=new XMLHttpRequest;a.open("GET",t,!1),a.setRequestHeader("Mlss-Userid",this.userId),a.onreadystatechange=function(e){if(4===this.readyState&&200===this.status){var t=JSON.parse(this.responseText).result;t.notebookAddress&&window.open(window.encodeURI(t.notebookAddress))}else{var a=JSON.parse(this.responseText);401===this.status||403===this.status?this.fetchError[this.status](a):this.$Message.error(a.message)}},a.send()},deleteItem:function(e){var t="/aide/".concat(this.FesEnv.aideApiVersion,"/namespaces/").concat(e.namespace,"/notebooks/").concat(e.name);this.deleteListItem(t)},copyJob:function(e){for(var t=this,a={},r=["cpu","namespace","queue","cluster","executorCores","executorMemory","executors","driverMemory"],o=0,s=r;o<s.length;o++){var n=s[o];a[n]=e[n]}a.memory=parseInt(e.memory),a.extraResources=e.gpu+"";var i=this.imageOptionList.indexOf(e.image)>-1?"Standard":"Custom";a.image={imageType:i,imageOption:"",imageInput:""},a.executorMemory=isNaN(parseInt(a.executorMemory))?"":parseInt(a.executorMemory)+"",a.driverMemory=isNaN(parseInt(a.driverMemory))?"":parseInt(a.driverMemory)+"","Standard"===i?a.image.imageOption=e.image:a.image.imageInput=e.image,setTimeout((function(){if(Object.assign(t.form,a),e.dataVolume){var r=t.dataVolumeInfo();for(var o in e.dataVolume)o>0&&t.addValidateItem(r.column,o,"dataVolume"),t.form.dataVolume[o]=e.dataVolume[o];t.storageDefalut="false",t.form.workspaceVolume.mountPath="/home/".concat(t.userId,"/workspace")}else t.form.workspaceVolume=e.workspaceVolume,t.storageDefalut="true"}),100),this.dialogVisible=!0},nextStep:function(){var e=this;this.$refs.formValidate.validate((function(t){t&&e.currentStep<e.stepList.length-1&&e.currentStep++}))},getListData:function(){var e=this;this.loading=!0;var t="",a={page:this.pagination.pageNumber,size:this.pagination.pageSize};if(this.query.namespace&&this.query.userName)t="/aide/".concat(this.FesEnv.aideApiVersion,"/namespaces/").concat(this.query.namespace,"/user/").concat(this.query.userName,"/notebooks");else if(this.query.namespace&&!this.query.userName)t="/aide/".concat(this.FesEnv.aideApiVersion,"/namespaces/").concat(this.query.namespace,"/notebooks");else{var r=localStorage.getItem("superAdmin");if("true"===r)t="/aide/".concat(this.FesEnv.aideApiVersion,"/namespaces/null/user/null/notebooks");else{var o=this.query.userName?this.query.userName:this.userId;t="/aide/".concat(this.FesEnv.aideApiVersion,"/user/").concat(o,"/notebooks")}}this.FesApi.fetch(t,a,"get").then((function(t){t&&"{}"!==JSON.stringify(t)&&(t.list&&t.list.forEach((function(e){var t=e.image.indexOf(":");e.image=e.image.substring(t+1),e.gpu=e.gpu||0,e.uptime=e.uptime?i["a"].transDate(e.uptime):"",e.queue&&(e.cluster=!0)})),e.dataList=Object.freeze(t.list),e.pagination.totalPage=t.total||0),e.loading=!1}),(function(){e.dataList=[],e.loading=!1}))},subInfo:function(){var e=this;this.$refs.formValidate.validate((function(t){if(t){e.setBtnDisabeld();var a={},r=["name","namespace"];r=e.form.cluster?r.concat(["queue","executorCores","executorMemory","executors","driverMemory"]):r;var o,i=Object(n["a"])(r);try{for(i.s();!(o=i.n()).done;){var l=o.value;("name"===l||"namespace"===l||e.form[l])&&(a[l]=e.form[l])}}catch(m){i.e(m)}finally{i.f()}a.cpu=parseFloat(e.form.cpu),a.memory={memoryAmount:parseInt(e.form.memory),memoryUnit:"Gi"},a.extraResources=JSON.stringify({"nvidia.com/gpu":e.form.extraResources});var c=e.form.image.imageOption||e.form.image.imageInput,u=e.defineImage+":"+c;a.image={imageName:u},a.image.imageType=e.form.image.imageType,"true"===e.storageDefalut?a.workspaceVolume=Object(s["a"])(Object(s["a"])({},e.form.workspaceVolume),{accessMode:"ReadWriteMany",mountType:"New"}):a.dataVolume=Object(s["a"])({},e.form.dataVolume[0]),e.FesApi.fetch("/aide/".concat(e.FesEnv.aideApiVersion,"/namespaces/namespace/notebooks"),a,"post").then((function(){e.getListData(),e.toast(),e.dialogVisible=!1}))}}))},dataVolumeInfo:function(){return{lg:this.form.dataVolume.length,column:["accessMode","localPath","subPath","mountPath","mountType"]}},changeImageType:function(e){"Custom"===e?this.form.image.imageOption="":this.form.image.imageInput=""},addValidateItem:function(e,t,a){var r,o=Object(n["a"])(e);try{for(o.s();!(r=o.n()).done;){var s=r.value;this.ruleValidate["".concat(a,"[").concat(t,"].").concat(s)]=this.ruleValidate["".concat(a,"[").concat(t-1,"].").concat(s)]}}catch(i){o.e(i)}finally{o.f()}},deleteValidateItem:function(e,t,a){var r,o=Object(n["a"])(e);try{for(o.s();!(r=o.n()).done;){var s=r.value;delete this.ruleValidate["".concat(a,"[").concat(t-1,"].").concat(s)]}}catch(i){o.e(i)}finally{o.f()}},clearDialogForm1:function(){var e=this,t=this.dataVolumeInfo();if(t.lg>1)for(var a=t.lg;a>1;a--)this.deleteValidateItem(t.column,a,"dataVolume");this.$refs.formValidate.resetFields();for(var r=["name","cpu","extraResources","memory","namespace","queue","driverMemory","executorCores","executorMemory","executors"],o=0,s=r;o<s.length;o++){var n=s[o];this.form[n]=""}for(var i in this.form.image)this.form.image[i]="imageType"===i?"Standard":"";this.form.workspaceVolume={accessMode:"",localPath:"",subPath:"",mountPath:"/home/".concat(this.userId,"/workspace"),mountType:"",size:0},this.form.dataVolume=[{accessMode:"",localPath:"",subPath:"",mountPath:"",mountType:"",size:0}],this.form.cluster=!1,setTimeout((function(){e.currentStep=0,e.storageDefalut="true",console.log("this.form",e.form)}),1e3)},editYarn:function(e){this.yarnDialog=!0,this.yarnForm={name:e.name,namespace:e.namespace,queue:e.queue||"",driverMemory:isNaN(parseInt(e.driverMemory))?"":parseInt(e.driverMemory)+"",executorCores:isNaN(parseInt(e.executorCores))?"":parseInt(e.executorCores)+"",executorMemory:isNaN(parseInt(e.executorMemory))?"":parseInt(e.executorMemory)+"",executors:isNaN(parseInt(e.executors))?"":parseInt(e.executors)+""}},saveYarnForm:function(){var e=this,t=this,a=Object(s["a"])({},this.yarnForm);a.driverMemory=Number(a.driverMemory)>0?a.driverMemory+"g":"",a.executorMemory=Number(a.executorMemory)>0?a.executorMemory+"g":"",this.$refs.yarnForm.validate((function(r){r&&e.FesApi.fetch("/aide/".concat(e.FesEnv.aideApiVersion,"/namespaces/").concat(t.yarnForm.namespace,"/notebooks"),a,{method:"patch",headers:{"Mlss-Userid":e.userId}}).then((function(){e.getListData(),e.toast(),e.yarnDialog=!1}))}))},cannelSaveYarnForm:function(){this.yarnDialog=!1,this.$refs.yarnForm.resetFormFields()}},beforeDestroy:function(){clearInterval(this.setIntervalList)}},w=E,P=(a("4d7f"),Object(p["a"])(w,r,o,!1,null,"2032f85e",null));t["default"]=P.exports}}]);
//# sourceMappingURL=chunk-1ce1c128.31b3e7b9.js.map