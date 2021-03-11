package main

import (
  "log"
  "math/rand"
  "os"
  "strconv"
  "strings"
)



type Neuron struct {
  tmin, thresholdMax, threshold int
  potential             int
  recovery, recoveryMax int
  inbox                 int
  fired bool
  targets []*Synapse
}

type Synapse struct {
  queue []byte
  signal int
  target *Neuron
  pointer int
}

func (s *Synapse) enqueue() {
  s.queue[s.pointer] = 1
}

func (s *Synapse) process() {
  s.pointer = (s.pointer + 1) % len(s.queue)
  if s.queue[s.pointer]>0 {
    s.target.enqueue(s.signal)
    s.queue[s.pointer] = 0
  }
}


func NewSynapse(n *Neuron, delay, signal int) *Synapse{
  return &Synapse{
    target: n,
    queue: make([]byte, delay),
    signal: signal,
  }
}

func NewNeuron(thresholdMax, recoveryMax int) *Neuron{
  return &Neuron{
    threshold:    1,
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

  if n.threshold > n.tmin{
    if n.recovery >= n.recoveryMax {
      n.threshold -=1
    } else {
      n.recovery +=1
    }
  }

}

func (n *Neuron) enqueue(signal int) {
  n.inbox += signal
}

func min(a, b int) int{
  if a < b {
    return a
  }
  return b
}

func main(){

  const NEURONS = 500
  const TIME = 5000


  var neurons []*Neuron
  var clefts []*Synapse

  var field [][]int

  for i := 0; i < NEURONS; i++ {
    neurons = append(neurons, NewNeuron(20, 16))
  }

  for _, n := range neurons{
    count := rand.Intn(11)+1
    for i := 0; i < count; i++ {
      target := neurons[rand.Intn(len(neurons))]
      cleft := NewSynapse(target, rand.Intn(7)+1, rand.Intn(5)-2)
      n.targets = append(n.targets, cleft)
      clefts = append(clefts, cleft)
    }
  }

  for i := 0; i < 10; i++ {
    neurons[i].enqueue(1)
  }

  for t := 0; t < TIME; t++ {

    current :=make([]int, NEURONS)
    field = append(field, current)
    for i, n := range neurons {
      n.process()
      if n.fired{
        current[i]=10
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
  for epoch, pots := range field {

    sb.Reset()
    sb.WriteString(strconv.Itoa(epoch))
    sb.WriteByte('\t')
    for _, pot := range pots {
      sb.WriteString(strconv.Itoa(pot))
      sb.WriteByte('\t')
    }
    sb.WriteByte('\n')
    f.WriteString(sb.String())
  }

  log.Println("Done")
}
