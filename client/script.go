package client

// This file is automatically generated! Do not modify it directly.

const Script = `
'use strict';var a = window.gogame = window.gogame || {};
var d = a.net = a.net || {};
d.Socket = function(c) {
  var b = this;
  b.b = new WebSocket(c);
  b.e = {};
  b.b.onmessage = function(e) {
    e = JSON.parse(e.data);
    if(e.i in b.e) {
      b.e[e.i](new f(e))
    }else {
      console.log(e)
    }
  };
  var h = "";
  b.b.onopen = function() {
    b.b.send(h);
    b.send = function(e) {
      b.b.send(JSON.stringify(e))
    }
  };
  b.send = b.send = function(b) {
    h += JSON.stringify(b) + "\n"
  };
  b.c = b.listen = function(e, c) {
    b.e[e] = c
  }
};
d.q = function() {
  return(d.g - 1).toString(32)
};
d.a = function() {
  d.g++;
  return d.q()
};
d.g = 0;
d.r = d.AttackerID = d.a();
d.o = d.VictimID = d.a();
d.h = d.Amount = d.a();
d.j = d.ChangeHealth = d.a();
d.s = d.a();
d.d = d.EntityID = d.a();
d.m = d.ParentID = d.a();
d.n = d.Tag = d.a();
d.l = d.EntitySpawned = d.a();
d.k = d.EntityDespawned = d.a();
d.f = d.EntityPosition = d.a();
d.FirstUnusedPacketID = d.a();
var f = d.Packet = function(c) {
  "object" == typeof c ? (this.i = c.i, this.p = c.p) : (this.i = c, this.p = {});
  this.set = this.set = function(b, c) {
    this.p[b] = c;
    return this
  };
  this.get = this.get = function(b) {
    return this.p[b]
  }
};
var g = a.client = {}, i = g.Entities = {};
g.disconnected = !1;
g.start = function(c) {
  g = a.client = new d.Socket(c);
  g.Entities = i;
  g.disconnected = !1;
  g.b.onerror = g.b.onclose = function() {
    g.disconnected = !0
  };
  g.c(d.l, function(b) {
    i[b.get(d.d)] = {parent:b.get(d.m), tag:b.get(d.n)}
  });
  g.c(d.k, function(b) {
    (function e(b) {
      delete i[b];
      for(var c in i) {
        i[c].parent == b && e(c)
      }
    })(b.get(d.d))
  });
  g.c(d.j, function(b) {
    i[b.get(d.o)].health = b.get(d.h)
  });
  g.c(d.f, function(b) {
    i[b.get(d.d)].position = b.get(d.f)
  })
};

`
