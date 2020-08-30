
var toString={}.toString,UA=window.navigator.userAgent.toLowerCase();

function isArray(v){
    return toString.call(v) == '[object Array]';
}

function isDomObject(obj){
    return  ( typeof HTMLElement === 'function' ) ? (obj instanceof HTMLElement):(obj && typeof obj === 'object' && obj.nodeType === 1 && typeof obj.nodeName === 'string');
}
function isUndefined(v){
    return toString.call(v) == '[object Undefined]'||(typeof v=="undefined");
}
function isObject(v){
    if(isUndefined(v)){
        return false;
    }
    return toString.call(v) == '[object Object]';
}
function isFunction(v){
    return toString.call(v) == '[object Function]';
}

function isChildDom(node,pnode){

    if(pnode.nodeType==1){
        if(node===pnode){
            return 1;
        }
        if(pnode.hasChildNodes){
            var sonnodes = pnode.childNodes;
            for (var i = 0; i < sonnodes.length; i++) {
                var sonnode = sonnodes.item(i);
                if(sonnode.nodeType == 1){
                    if(sonnode===node){
                        return true;
                    }
                    var rt=isChildDom(node,sonnode);
                    if(rt){
                        return rt;
                    }
                }
            }
        }
    }

}
function bubbleFind(node,each,pnode,proxy){

    //冒泡查找
    if(!(isDomObject(node)&&isFunction(each))){
        return false;
    }
    pnode=pnode||document.documentElement;
    proxy=proxy||window;

    var target=null;

    while(node!==null&&node!=pnode){

        if(node.nodeType==1){

            if(each.call(proxy,node)){
                target=node;
                break;
            }
        }
        node=node.parentNode;
    }

    return target;

}


function forIn(obj,each,proxy){
    proxy=proxy||obj;
    var i,rt;
    if(isObject(obj)){
        for(i in obj){
            if(obj.hasOwnProperty(i)){
                if(isFunction(each)){
                    rt=each.call(proxy,i,obj[i]);
                    if(rt=="__break"){
                        break;
                    }
                    if(rt=="__continue"){
                        continue;
                    }
                }
            }
        }
        return;
    }
    if(isArray(obj)){
        for(i=0;i<obj.length;i++){
            if(isFunction(each)){
                rt=each.call(proxy,i,obj[i]);
                if(rt=="__break"){
                    break;
                }
                if(rt=="__continue"){
                    continue;
                }
            }
        }
    }

}
function Copy(obj){
    if(typeof(obj)!="object" || obj===null)return obj;
    var o={};
    if(isArray(obj)){
        o=[];
    }
    forIn(obj,function(i,v){
        if(isArray(v)||(isObject(v)&& v!==null&&!isDomObject(v))){
            o[i]=Copy(v);
        }else{
            o[i]=v;
        }
    });
    return o;

}

function Xl(p){

    //访问Xl属性，如果不存，则到配置文件中的搜寻
    if(Xl.isEmpty(p)){return Xl;} //返回自身

}
Xl.extend=function(target, source){
    var stype=typeof source;
    if (typeof target!=="object"&&!isFunction(target) ) {
        target = {};
    }
    if(stype==='undefined'||stype ==='boolean'){
        source=target;
        target=this;
    }
    var cs=Copy(source);
    for (var p in cs) {
        if (cs.hasOwnProperty(p)) {
            target[p] = cs[p];
        }
    }
    return target;
};

Xl.inherit=function(target,parent){

    parent=Copy(parent);
    for (var p in parent){
        if (parent.hasOwnProperty(p)) {
            if(Xl.isUndefined(target[p])){
                target[p] = parent[p];
            }
        }
    }
    target.parent=parent;

    return target;
};

Xl.extend({
    VERSION:'2.0.0',
    isHtm5:window.applicationCache?true:false,
    Ready:{},
    Cache: {},
    Promise:window.Promise,
    forIn:forIn,
    copy:Copy,
    getBrowerV:function(){
        if(!Xl.isEmpty(Xl._browerVersion)){return Xl._browerVersion;}
        var b=Xl._browerVersion={},u= UA;
        var s;
        (s = u.match(/msie ([\d.]+)/)) ? (b.ie = s[1]) :
            (s = u.match(/firefox\/([\d.]+)/)) ? (b.firefox = s[1]) :
                (s = u.match(/chrome\/([\d.]+)/)) ? (b.chrome = s[1]) :
                    (s = u.match(/opera.([\d.]+)/)) ? (b.opera = s[1]) :
                        (s = u.match(/version\/([\d.]+).*safari/)) ? (b.safari) = s[1] : 0;
        return b;
    },
    isIE8:function(){
        //是否是ie8以下版本包含ie8
        var v=Xl.getBrowerV();
        return v.ie?(parseInt(v.ie)<9?true:false):false;
    },
    getGuid:function(){
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g,function(c){
            var r=Math.random()*16|0;
            var v=c=='x'?r:(r&0x3|0x8);
            return v.toString();
        }).toUpperCase();
    },
    setG:function(key, value) {
        function P(r, k, v) {
            var ka = k.split(">");
            var len = ka.length;
            if (len == 1) {
                r[ka[0]] = v;
            } else if (len == 2) {
                r[ka[0]] ? '' : (r[ka[0]] = {});
                r[ka[0]][ka[1]] = v;
            }
        }
        if (typeof key == "string") {
            var K = key.split('/');
            var len = K.length;
            switch (len) {
                case 1:
                    P(Xl.Cache, K[0], value);
                    break;
                case 2:
                    Xl.Cache[K[0]] ? '' : (Xl.Cache[K[0]] = {});
                    P(Xl.Cache[K[0]], K[1], value);
                    break;
                case 3:
                    Xl.Cache[K[0]] ? '' : (Xl.Cache[K[0]] = {});
                    Xl.Cache[K[0]][K[1]] ? '' : (Xl.Cache[K[0]][K[1]] = {});
                    P(Xl.Cache[K[0]][K[1]], K[2], value);
                    break;
            }
        }
    },
    getG:function(key) {
        if (typeof key == "string") {
            var K = key.split('/');
            var len = K.length;
            try {
                switch (len) {
                    case 1:
                        return Xl.Cache[K[0]];
                    case 2:
                        return Xl.Cache[K[0]][K[1]];
                    case 3:
                        return Xl.Cache[K[0]][K[1]][K[2]];
                }
            } catch (err) {
                return 'undefined';
            }
        }
    },
    _getGArr:function(key){
        var arr=Xl.getG(key);
        if(Xl.isEmpty(arr)||!Xl.isArray(arr)){
            arr=[];
        }
        return arr;
    },
    pushG:function(key,v){
        var vs=Xl._getGArr(key);
        vs.push(v);
        Xl.setG(key,vs);
        return vs;
    },
    popG:function(key){
        var vs=Xl._getGArr(key);
        vs.pop();
        Xl.setG(key,vs);
        return vs;
    },
    shiftG:function(key){
        var vs=Xl._getGArr(key);
        vs.shift();
        Xl.setG(key,vs);
        return vs;
    },
    unshiftG:function(key,v){
        var vs=Xl._getGArr(key);
        vs.unshift(v);
        Xl.setG(key,vs);
        return vs;
    },
    isUndefined:isUndefined,
    isBoolean:function(v){
        return toString.call(v) == '[object Boolean]';
    },
    isNumber:function(v){
        if(Xl.isString(v)){
            if(/^\s*((-|\+)?)\d+?(\.?)\d*\s*$/.test(v)){
                return true;
            }
        }
        return toString.call(v) == '[object Number]';
    },
    isString:function(v){
        return toString.call(v) == '[object String]';
    },
    isArray:isArray,
    isFunction:isFunction,
    isObject:isObject,
    isRegExp:function(v){
        return toString.call(v) == '[object RegExp]';
    },
    isPlainObject:function(v){
        if(!Xl.isObject(v)){
            return false;
        }
        for(var i in v){
            if(Xl.isFunction(v[i])){
                return false;
            }
        }
        return true;
    },
    isDomObject:isDomObject,
    isChildDom:isChildDom,
    bubbleFind:bubbleFind,
    inArray:function(value,array,ct){
        var rt=-1;
        Xl.forIn(array,function(i,v){
            if(ct){
                if(v===value){
                    rt=i;
                }
            }else{
                if(v==value){
                    rt=i;
                }
            }
        });
        if(rt==-1){
            return false;
        }
        return true;
    },
    isEmpty:function(at){
        if(typeof at=="undefined"){return true;}
        if(at===null){return true;}
        if(typeof at=="string"){
            if(/^\s*$/g.test(at)){return true;}
        }
        if(Xl.isArray(at)){if(at.length===0){return true;}}
        if(Xl.isPlainObject(at)){for(var i in at){return false;}
            return true;}
    },
    leftStr:function(a,len,b){

        if(!Xl.isString(a)){
            return a;
        }
        if(Xl.isUndefined(b)){
            b=true;
        }
        if(b){
            return a.substr(0,len);
        }else if(a.length>len){
            return a.substr(0,len)+'...';
        }else{
            return a;
        }
    },
    trim:function(s){
        if(!Xl.isString(s)){return s;}
        return s.replace(/(^\s*)|(\s*$)/g,'');
    },
    removeFrom:function(obj,find){

        if(Xl.isArray(obj)){
            obj=obj.filter(function(x){return x!=find;});
        }
        return obj;

    },
    addDivToBody:function(id, tag) {
        tag = tag || 'div';
        var oBody = document.getElementsByTagName('BODY').item(0);
        var odiv = document.createElement(tag);
        odiv.setAttribute('id', id);
        odiv.guid = Xl.getGuid();
        oBody.appendChild(odiv);
        return odiv;
    },
    fetchObjByKeys:function(keys,obj){
        if(Xl.isString(keys)){
            keys=keys.split(',');
        }
        var tem={};
        if(Xl.isArray(keys)){
            for(var k in keys){
                if(obj.hasOwnProperty&&obj.hasOwnProperty(keys[k])){
                    tem[keys[k]]=obj[keys[k]];
                }else{
                    tem[keys[k]]=obj[keys[k]];
                }
            }
        }
        return tem;
    },
    accessObjFromKeys:function(obj,arr){

        if(!Xl.isEmpty(arr)){
            arr.forEach(function(p){
                if(obj.hasOwnProperty(p)){
                    obj=obj[p];
                }else{
                    Xl.alert("访问的属性"+p+"不存在");
                }
            });
        }
        return obj;
    },
    getViewSize:function(obj) {
        obj = Xl.E(obj)||document.documentElement;
        var scrollTop,scrollLeft;
        if (obj == document.documentElement) {
            scrollTop = obj.scrollTop||document.body.scrollTop;
            scrollLeft = obj.scrollLeft||document.body.scrollLeft;
        }
        var screen = window.screen;
        var vsize=Xl.fetchObjByKeys("clientWidth,clientHeight,offsetWidth,offsetHeight",obj);
        vsize.screen=screen;
        vsize.scrollTop=scrollTop;
        vsize.scrollLeft=scrollLeft;

        return vsize;
    },
    centerWindow:function(ob, w, h) {
        var sn = Xl.getViewSize();
        var bh = (sn.clientHeight - h) / 2;
        bh = bh < 0 ? 0 : bh;
        var scrolltop = sn.scrollTop + bh;
        var bw = Math.abs((sn.clientWidth - w) / 2);
        ob.style.left = bw + "px";
        ob.style.top = scrolltop + "px";
    },
    date:function(fmt,time){
        if(Xl.isUndefined(time)){
            time=Xl.getTime();
        }else{
            time=time*1000; //转换为毫秒
        }
        if(Xl.isUndefined(fmt)){
            fmt="yyyy-MM-dd hh:mm:ss";
        }
        var t=new Date(time);
        var o = {
            "M+" : t.getMonth()+1,
            "d+" : t.getDate(),
            "h+" : t.getHours(),
            "m+" : t.getMinutes(),
            "s+" : t.getSeconds(),
            "q+" : Math.floor((t.getMonth()+3)/3),
            "S"  : t.getMilliseconds()
        };
        if(/(y+)/.test(fmt)){
            fmt=fmt.replace(RegExp.$1, (t.getFullYear()+"").substr(4 - RegExp.$1.length));
        }
        for(var k in o){
            if(new RegExp("("+ k +")").test(fmt)){
                fmt = fmt.replace(RegExp.$1, (RegExp.$1.length==1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));
            }
        }
        return fmt;
    },
    getTime:function(){
        return (new Date()).getTime();
    },
    strToTime: function(date) {
        if (date === 0) {
            return 0;
        }
        var moth,d;
        if (date) {
            var datearr = date.split('-');
            if (datearr[1].substr(0, 1) == "0") {
                moth = parseInt(datearr[1].substr(1, 1)) - 1;
            } else {
                moth = parseInt(datearr[1]) - 1;
            }
            d = new Date(datearr[0], moth, datearr[2]);
            return d.getTime();
        } else {
            d = new Date();
            var year = d.getFullYear();
            moth = d.getMonth();
            var day = d.getDate();
            d = new Date(year, moth, day);
            return d.getTime();
        }
    },
    getDJSformat: function(endtime) {
        if (!/^\d+$/g.test(endtime)) {
            endtime = Xl.strToTime(endtime);
        } else {
            endtime = endtime * 1000;
        }
        var d = new Date();
        var nowtime = d.getTime();
        if (nowtime > endtime) {
            return {timeend:true};
        }
        var cm = (endtime - nowtime) / 1000;
        var day = Math.floor(cm / (3600 * 24));
        if (day > 365 * 10) {
            return {timeend:false,day:day,hour:0,minute:0,second:0};
        }
        cm = cm - day * 3600 * 24;
        var hour = Math.floor(cm / 3600);
        cm = cm - hour * 3600;
        var minute = Math.floor(cm / 60);
        cm = cm - minute * 60;
        var second = parseInt(cm);
        return {timeend:false,day:day,hour:hour,minute:minute,second:second};
    },
    getImageSize: function(src, func) {
        var i = new Image();
        i.src = src;
        if (Xl.isFunction(func)) {
            if (i.complete) {
                func({
                    w: i.width,
                    h: i.height,
                    src: src
                });
            } else {
                i.onload = function() {
                    func({
                        w: i.width,
                        h: i.height,
                        src: src
                    });
                };
            }
        } else {
            return {
                w: i.width,
                h: i.height,
                src: src
            };
        }
    },
    zoomSize:function(sw,sh,boxw,boxh){
        var w,h;
        boxw=boxw||20000;boxh=boxh||20000;
        if(sw<=boxw&&sh<boxh){return {w:sw,h:sh};}
        var b1=boxw/sw;var b2=boxh/sh;
        var b=sw/sh;
        if(b1<=b2){
            w=boxw;h=w/b;return{w:w,h:h};
        }else{
            h=boxh;w=b*h;
            return{w:w,h:h};
        }
    },
    sgData:function(id,d,v){
        var dom=null;
        dom=Xl.E(id);
        if(!dom){return '';}
        if(dom.dataset){
            if(Xl.isUndefined(v)){
                return dom.dataset[d];
            }else{
                dom.dataset[d]=v;
            }
        }else{
            if(Xl.isUndefined(v)){
                return dom.getAttribute("data-"+d);
            }else{
                dom.setAttribute("data-"+d,v);
            }
        }
    },
    getLen:function(obj){
        if(Xl.isNumber(obj)){
            obj=obj.toString();
        }
        if(Xl.isArray(obj)||Xl.isString(obj)){
            return obj.length;
        }else if(Xl.isPlainObject(obj)){
            var len=0;
            for(var i in obj){
                if(obj.hasOwnProperty(i)){
                    len++;
                }
            }
            return len;
        }
        return 0;
    },
    b64EncodeUnicode(str) {
        return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, 
                          function(match, p1) {
                               return String.fromCharCode('0x' + p1);
                          }));
   }

});

Xl.Mem={
    set:function(key,value){
        var keyarr=[];
        if(key.indexOf('/')){
            keyarr=key.split('/');
        }else{
            keyarr=[key];
        }
        var ov=Xl.store.get(keyarr[0])||{};
        var keylen=Xl.getLen(keyarr);
        switch(keylen){
            case 2:
                ov[keyarr[1]]=value;
                break;
            case 3:
                ov[keyarr[1]]=ov[keyarr[1]]||{};
                ov[keyarr[1]][keyarr[2]]=value;
                break;
            default:
                ov=value;
                break;
        }
        Xl.store.set(keyarr[0],ov);
    },
    get:function(key){
        var keyarr=[];
        if(key.indexOf('/')){
            keyarr=key.split('/');
        }else{
            keyarr=[key];
        }
        var ov=Xl.store.get(keyarr[0])||{};
        var keylen=Xl.getLen(keyarr);
        switch(keylen){
            case 2:
                return ov[keyarr[1]];
            case 3:
                ov[keyarr[1]]=ov[keyarr[1]]||{};
                return ov[keyarr[1]][keyarr[2]];
            default:
                return ov;
        }
    }

};

Xl.Reg={
    isEmail:function(str){
        return /\s*[a-zA-Z0-9]+@[a-z0-9]+?\.[a-z0-9]\s*/g.test(str);
    },
    isTel:function(str){
        return /^\s*1(3|4|5|6|7|8|9)\d{9}\s*$/g.test(str);
    },
    isHanzi:function(str,ispart){
        if(ispart){
            return /[\u4e00-\u9fa5]/g.test(str);
        }
        return /^\s*[\u4e00-\u9fa5]+\s*$/g.test(str);
    }
};


export default Xl;
