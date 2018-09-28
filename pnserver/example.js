<p>Enter a number and click OK:</p>

<form> 

<input id="id1" type="number" max="100">
<button onclick="myFunction()">OK</button>
<select id="xx" onchange="MyFunc2()">
<option name="A" value="A">A</option>
<option name="B" value="B">B</option>
</select>

<select id="yy">
<option name="X" value="X">X</option>
<option name="Y" value="Y">Y</option>
</select>


</form>


<p>If the number is greater than 100 (the input's max attribute), an error message will be displayed.</p>

<p id="demo"></p>

<p id="msg2"></p>
<p id="msg3"></p>
<p id="msg4"></p>

<script>
function myFunction() {
    var txt = "";
    if (document.getElementById("id1").validity.rangeOverflow) {
        txt = "Value too large";
    } else {
        txt = "Input OK";
    } 
    document.getElementById("demo").innerHTML = txt;
}

function ClearSelBox(sel) {
    var n = sel.length;
    var i;
    for (i = 0; i < n; i++) {
        sel.remove(0);
    }
}

function FillSelBox(sel, lst) {
    var i;
    var n = lst.length;
    for (i = 0; i < n; i++) {
        var option = document.createElement("option");
        option.text = lst[i];
        sel.add(option);
    }
}

function MyFunc2() {
    var txt = "";
    var sel = document.getElementById("xx").value;
    var options;
    if (sel == "A") {
      options = ["Dog", "Cat", "Mouse"];
     } else {
      options = ["Green", "Blue", "Red"];
     }
    sel2 = document.getElementById("yy");
    ClearSelBox(sel2);
    document.getElementById("msg3").innerHTML = "Calling FillSelBox";
    FillSelBox(sel2, options);
    document.getElementById("msg2").innerHTML = "value changed: " + sel;
    }

</script>

</body>
</html>






<!DOCTYPE html>
<html>
<body>

<p id="outmsg">[nothing]</p>

<script type="application/json" id="projects">
    {
        "C19 -- Competition 2018": ["CH", "AH", "TW", "AM", "EL"],
        "P18 -- Prototype 2018": ["CH", "XY", "BQ", "R8"],
        "F18 -- Fall 2018": ["CH -- Chassis", "FS -- Fresby Shooter", "MS -- Misc"]
    }
</script>

<script type="text/javascript">

function clearSelBox(sel) {
    var n = sel.length;
    var i;
    for (i = 0; i < n; i++) {
        sel.remove(0);
    }
}

function fillSelBox(sel, lst) {
    var i;
    var n = lst.length;
    for (i = 0; i < n; i++) {
        var option = document.createElement("option");
        option.text = lst[i];
        sel.add(option);
    }
}

function fillProjects() {
    var data = JSON.parse(document.getElementById("projects").innerHTML);
    var prjlst = Object.keys(data);
    var sel = document.getElementById("nep_project_field");
    clearSelBox(sel);
    fillSelBox(sel, prjlst);
    document.getElementById("outmsg").InnerHTML = "Got here 2";
}

function fillSubsystem() {
    var data = JSON.parse(document.getElementById("projects").innerHTML);
    var p = document.getElementById("nep_project_field");
    var sel = document.getElementById("nep_subsystem_field");
    clearSelBox(sel);
    fillSelBox(sel, data[p]);
}
</script>

<p>Here we go. </p>

<select id="nep_project_field" name="Project" onchange="fillSubsystem()">  </select>
<br>
<br>
<select id="nep_subsystem_field" name="Subsystem"></select>


<script type="text/javascript">
fillProjects();
document.getElementById("outmsg").InnerHTML = "Got here";
fillSubsystem();
</script>

</body>
</html> 