package main

import (
  "flag"
  "log"
  "math/rand"
  "os"
  "strconv"
  "strings"
    "fmt"
)



type Neuron struct {
  thresholdMin, thresholdMax, threshold int
  potential             int
  recovery, recoveryMax int
  inbox                 int
  fired bool
  targets []*Synapse
}

type Synapse struct {
  queue uint
  signal int8
  target *Neuron
  pointer int8
  size int8
}

func (s *Synapse) enqueue() {
  s.queue |= 1<<s.pointer
}

func (s *Synapse) process() {
  // If synapse is empty, then we can skip processing entirely
  if s.queue == 0{
    return
  }

  s.pointer+=1
  if s.pointer == s.size{
    s.pointer = 0
  }

  mask := uint(1<<s.pointer)
  if s.queue & mask != 0 {
    s.queue &= ^mask
    s.target.enqueue(s.signal)
  }
}


func NewSynapse(n *Neuron, delay int8, signal int8) *Synapse{
  return &Synapse{
    target: n,
    size:delay,
    signal: signal,
  }
}

func NewNeuron(thresholdMax, recoveryMax, thresholdMin int) *Neuron{
    
  return &Neuron{    
    threshold:    thresholdMin,
    thresholdMin: thresholdMin,
    recoveryMax:  recoveryMax,
    thresholdMax: thresholdMax,
  }
}

const DROP = -2

func (n *Neuron) process() {
  signal := n.inbox
  n.inbox = 0

  if signal > 0 && n.potential >= 0 {
    n.potential += signal
  } else if n.potential > 0 {
    n.potential -= 1
  } else if n.potential < 0 {
    n.potential += 1
  }


  if n.potential >= n.threshold {
    n.potential = DROP
    n.threshold = min(n.thresholdMax, n.threshold+1)
    n.recovery = 0
    n.fired = true

    for _, s := range n.targets{
      s.enqueue()
    }

  } else {
    n.fired = false
  }

  if n.threshold > n.thresholdMin {
    if n.recovery >= n.recoveryMax {
      n.threshold -=1
      n.recovery = 0
    } else {
      n.recovery +=1
    }
  }

}

func (n *Neuron) enqueue(signal int8) {
  n.inbox += int(signal)
}

func min(a, b int) int{
  if a < b {
    return a
  }
  return b
}


var(
  neuronCount = flag.Int("neurons", 4000, "neuron count")
  time = flag.Int("time", 10000, "max ticks")


  recMultiplier = flag.Int("rec", 10, "recovery multiplier")
  tmaxMultiplier = flag.Int("tmax", 10, "tmax multiplier")
  tminBound = flag.Int("tmin", 3, "tmin upper bound")
)

func main(){

  flag.Parse()
  if !flag.Parsed(){
    log.Panicln("problem")
  }

   NEURONS := *neuronCount
  TIME := *time

  CLUSTER := NEURONS/8


  var neurons []*Neuron
  var clefts []*Synapse

  var field [][]int

    
  for i := 0; i < NEURONS; i++ {
    cluster := i / CLUSTER
      
    tmin := rand.Intn(*tminBound)
    tmax := tmin + cluster* *tmaxMultiplier
    rec := 3+cluster* *recMultiplier
      
    neurons = append(neurons, NewNeuron(tmax, rec, tmin))
  }

  for j, n := range neurons{
    // local connectivity
    cluster := j / CLUSTER
    
    
    for i := 0; i < 20; i++ {

      lid := rand.Intn(CLUSTER) + cluster * CLUSTER
      if lid != j { 
        target := neurons[lid]
        cleft := NewSynapse(target, int8(rand.Intn(3+cluster/2)+1), int8(rand.Intn(5)-2))
        n.targets = append(n.targets, cleft)
        clefts = append(clefts, cleft)
      }
    }
    
      
    
    // inter-cluster
    // higher thoughts have bigger connectivity
    connections := cluster * 2
    for i := 0; i < connections; i++ {
      target := neurons[rand.Intn(NEURONS)]
      cleft := NewSynapse(target, int8(rand.Intn(4)+cluster+1), int8(rand.Intn(4)-1))
      n.targets = append(n.targets, cleft)
      clefts = append(clefts, cleft)
    }
    
    // forward connectivity
    for i:= 0; i < 4; i++ {
      mx := NEURONS-j
      lid := rand.Intn(mx+1) + j
        target := neurons[lid % NEURONS]
        cleft := NewSynapse(target, int8(rand.Intn(3+cluster/2)+1), int8(rand.Intn(4)-1))
        n.targets = append(n.targets, cleft)
        clefts = append(clefts, cleft)
    }
  }
    


  for i := 0; i < 10; i++ {
    neurons[i].enqueue(1)
  }
  
  pmax := 30

  for t := 0; t < TIME; t++ {

    current :=make([]int, NEURONS)
    field = append(field, current)
    for i, n := range neurons {
      n.process()
      if n.potential > pmax{
        pmax = n.potential
      }
      if n.fired{
        current[i]=pmax+5
      } else {
        current[i]=n.potential
      }
    }
    for _, c := range clefts {
      c.process()
    }
  }


  f, err := os.Create("dump.tsv")
  if err != nil {
    log.Panicln(err)
  }
  defer f.Close()

  var sb strings.Builder
  for _, pots := range field {

    sb.Reset()
    for r, pot := range pots {
        if r != 0 {
      sb.WriteByte('\t')      
        }
      sb.WriteString(strconv.Itoa(pot))
          }
    sb.WriteByte('\n')
    f.WriteString(sb.String())
  }

    fmt.Println(fmt.Sprintf("max: %v, clefts: %d, neurons: %d, ratio: %d", pmax, len(clefts), NEURONS, len(clefts)/NEURONS))
}
