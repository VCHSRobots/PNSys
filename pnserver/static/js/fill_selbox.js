// --------------------------------------------------------------------
// fill_selbox.js -- Fill various selection boxes
//
// Created 2018-10-01 DLB
// --------------------------------------------------------------------

function clearSelBox(sel) {
    var n = sel.length;
    var i;
    for (i = 0; i < n; i++) {
        sel.remove(0);
    }
}

function addselection(selbox, txt, lab, val) {
    var option = document.createElement("option");
    option.text = txt;
    option.label = lab;
    option.value = val;
    selbox.add(option);
}

function countMatchLen(s1, s2) {
  n = s1.length;
  if (n > s2.length) { n = s2.length;}
  for(var i = 0; i < n; i++) {
    c1 = s1.substr(i, 1).toLowerCase();
    c2 = s2.substr(i, 1).toLowerCase();
    if (c1 != c2) {
      return i;
    }
  }
  return n;
}

function setSelectionBox(selid, txt) {
    //console.log("in setSelectioBox. Txt=", txt)
    var selbox = document.getElementById(selid);
    var i;
    // Look for an exact match, in either text or values.
    for(i = 0; i < selbox.options.length; i++) {
      var opt = selbox.options[i];
      if (opt.value == txt) {
          selbox.value = opt.value;
          //console.log("setting to (1):", selbox.value)
          return true;
      }
      if (opt.text == txt) {
          selbox.value = opt.value;
          //console.log("setting to (2):", selbox.value)
          return true;
      }
    }
    if (txt == "") {
        // This means that a correct default was not 
        // provided, so the current selection is valid.
        return true
    }
    // Return false if no valid selection was found.
    return false;
}

function setSelectionBoxBestMatch(selid, txt) {
    //console.log("setting ", selid, "  with ", txt)
    if (setSelectionBox(selid, txt)) {
        //console.log("set with exact match.");
        return true;
    }
    // Try again but by finding the closest match.
    var bestopt = false;
    var nchars = -1;
    for (i = 0; i < selbox.options.length; i++) {
      var opt = selbox.options[i];
      var matchlen = countMatchLen(txt, opt.text);
      if (matchlen > nchars) {
        nchars = matchlen;
        bestopt = opt;
      }
    }
    if (bestopt && nchars > 0) {
      selbox.value = bestopt.value;
      //console.log("setting to (3):", selbox.value, "  Matched chars", nchars);
      return true;
    }
    //console.log("unable to set selection box?");
    return false;
}

function setTextBox(fieldid, txt) {
    var d = document.getElementById(fieldid);
    d.text = txt;
    d.value = txt ;
}

// fillDesigers uses the JSON from "designers" to fill a selection box 
function fillDesigners(selbox_id, showall, add_blank=false) {
    dobj = document.getElementById("designers");
    if (typeof dobj === "undefined") {
        console.log("fillDesigners, designers json unknown.")
        return;
    }
    var selbox = document.getElementById(selbox_id);
    if (typeof selbox === "undefined") {
        console.log("fillDesigners, selbox id unknown: ", selbox_id)
        return;
    }
    var desingerlst = JSON.parse(dobj.innerHTML);
    clearSelBox(selbox);
    if (add_blank) { addselection(selbox, "", "", ""); }
    var i;
    for(i = 0; i < desingerlst.length; i++) {
        d = desingerlst[i];
        if (showall || d.Active) {
            addselection(selbox, d.Name, d.Name, d.Name);
        }
    }
}

// fillCategories uses the JSON from "categories" to fill a selection box.
function fillCategories(selbox_id, add_blank=false) {
    cobj = document.getElementById("categories");
    if (typeof cobj === "undefined") {
        console.log("fillCategories, categories json unknown,")
        return
    }
    var selbox = document.getElementById(selbox_id);
    if (typeof selbox === "undefined") {
        console.log("fillCategories, selbox id unknown: ", selbox_id)
        return;
    }
    var catlst = JSON.parse(cobj.innerHTML);
    clearSelBox(selbox);
    if (add_blank) { addselection(selbox, "", "", ""); }
    var i;
    for (i = 0; i < catlst.length; i++) {
        var label = catlst[i].Category + " -- " + catlst[i].Description;
        addselection(selbox, label, label, catlst[i].Category);
    }
}


function fillPartTypes(selbox_id, add_blank=false) {
    pobj = document.getElementById("parttypes")
    if (typeof pobj === "undefined") {
        console.log("fillPartTypes, partypes json unknown.")
        return
    }
    var selbox = document.getElementById(selbox_id);
    if (typeof selbox === "undefined") {
        console.log("fillPartTypes, selbox id unknown: ", selbox_id)
        return;
    }
    var ptlst = JSON.parse(pobj.innerHTML);
    clearSelBox(selbox);
    if (add_blank) { addselection(selbox, "", "", ""); }
    var i;
    for (i = 0; i < ptlst.length; i++) {
        var option = document.createElement("option");
        var label = ptlst[i].Digit + " -- " + ptlst[i].Description;
        addselection(selbox, label, label, ptlst[i].Digit);
    }
}

function fillVendors(selbox_id, add_blank=false, add_newvendor=false) {
    vobj = document.getElementById("vendors")
    if (typeof vobj === "undefined") {
        console.log("fillVendors, vendors json unknown.")
        return
    }
    var selbox = document.getElementById(selbox_id);
    if (typeof selbox === "undefined") {
        console.log("fillVendors, selbox id unknown: ", selbox_id)
        return;
    }
    var vendorlst = JSON.parse(vobj.innerHTML);
    clearSelBox(selbox);
    if (add_blank) { addselection(selbox, "", "", ""); }
    var i;
    for(i = 0; i < vendorlst.length; i++) {
        vendor = vendorlst[i];
        addselection(selbox, vendor, vendor, vendor);
    }
    if (add_newvendor) {
        addselection(selbox, "<new vendor>", "<new vendor>", "<new vendor>")
    }
}

function fillSubsystem(selbox_id, project_box_id, add_blank=false) {
    pobj = document.getElementById("projects")
    if (typeof pobj === "undefined") {
        console.log("fillSubsystem, projects json unknown.")
        return
    }
    var prjlst = JSON.parse(pobj.innerHTML);
    var selbox = document.getElementById(selbox_id);
    if (typeof selbox === "undefined") {
        console.log("fillSubsystem, selbox id unknown: ", selbox_id)
        return;
    }
    var prjbox = document.getElementById(project_box_id);
    if (typeof prjbox === "undefined") {
        console.log("fillSubsystem, project_box_id unknown: ", project_box_id)
        return;
    }
    clearSelBox(selbox);
    if (add_blank) { addselection(selbox, "", "", ""); }
    var i;
    for (i = 0; i < prjlst.length; i++) {
        if (prjlst[i].ProjectId == prjbox.value) {
            var j; 
            var subs = prjlst[i].Subsystems;
            for (j = 0; j < subs.length; j++) {
                var label = subs[j].SubsystemId + " -- " + subs[j].Description;
                addselection(selbox, label, label, subs[j].SubsystemId);
            }
            return;
        }
    }
}

function fillProjects(selbox_id, showall, add_blank=false) {
    pobj = document.getElementById("projects")
    if (typeof pobj === "undefined") {
        console.log("fillProjects, projects json unknown.")
        return
    }
    var prjlst = JSON.parse(pobj.innerHTML);
    var selbox = document.getElementById(selbox_id);
    if (typeof selbox === "undefined") {
        console.log("fillProjects, selbox id unknown: ", selbox_id)
        return;
    }
    clearSelBox(selbox);
    if (add_blank) { addselection(selbox, "", "", ""); }
    var i;
    for(i = 0; i < prjlst.length; i++) {
        prj = prjlst[i];
        if (showall || prj.Active) {
            var label = prj.ProjectId + " -- " + prj.Description;
            addselection(selbox, label, label, prj.ProjectId);
        }
    }
}
