/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

function demo_square()

    my_handle             = scf(100069);
    clf(my_handle,"reset");



    my_plot_desc          = "ModbusSim signal Square";
    my_handle.figure_name = my_plot_desc;

//Tomaremos la muestra por punto de muestro
//En nuestro caso se muestreara a 100 ms = 10 Hz
    frecuencia = 0.5 //Hz
    fase = 1;
    paso = 0;
    t=[];
    tick=[];    
    for x=1:1:100
     t(x) = paso;  
     tick(x)=x;         
     paso = paso + 2*%pi*frecuencia/10;  
     if (paso>=2*%pi) then
         paso=paso - 2*%pi;;  
     end;           
    end
    
    f = sign(sin(t + fase));
    
    plot2d(tick',f');
    xtitle(my_plot_desc,"t","square(t+fase)");
    xgrid(color("grey"));

endfunction

demo_square();
clear demo_si;
